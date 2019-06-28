package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	actorlog "github.com/AsynkronIT/protoactor-go/log"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/worker"
)

func main() {
	os.Exit(realMain())
}

var (
	masterHost    = flag.String("host", "", "host address and port of master worker")
	timeoutSec    = flag.Int("timeout", 5, "timeout (seconds)")
	watchDuration = flag.Int("duration", 300, "ping duration (milliseconds)")
	listenAddr    = flag.String("listen", "127.0.0.1:8888", "listen address")
	degub         = flag.Bool("debug", false, "debug mode")
	// TODO: maxStep         = flag.Int("max-step", 0, "maximum number of supper step, 0 means no limit")
)

func realMain() int {
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()

	if *degub {
		actor.SetLogLevel(actorlog.DebugLevel)
		remote.SetLogLevel(actorlog.DebugLevel)
	} else {
		actor.SetLogLevel(actorlog.ErrorLevel)
		remote.SetLogLevel(actorlog.ErrorLevel)
	}

	remote.Start(*listenAddr)

	// parse flags
	if *masterHost == "" {
		flag.PrintDefaults()
		exitErr(errors.New("no host specified"))
	}

	// common setup
	coordinator := actor.NewPID(*masterHost, worker.CoordinatorActorID)

	// run command
	var err error
	switch {
	case len(args) == 0:
		err = errors.New("no command specified")
	case args[0] == "state":
		err = showStat(coordinator)
	case args[0] == "load":
		load(coordinator, args[1:])
	case args[0] == "start":
		err = startSuperStep(coordinator)
	case args[0] == "watch":
		err = watch(coordinator)
	case args[0] == "agg":
		err = showAggregatedValue(coordinator)
	case args[0] == "shutdown":
		err = sendShutdown(coordinator)
	// TODO: help command
	default:
		err = fmt.Errorf("%s - no such command", args[0])
		flag.PrintDefaults()
	}

	if err != nil {
		exitErr(err)
	}

	return 0
}

func exitErr(err error) {
	log.Fatalf("ERROR: %v", err)
}

func showAggregatedValue(coordinator *actor.PID) error {
	ctx := actor.EmptyRootContext
	fut := ctx.RequestFuture(coordinator, &command.ShowAggregatedValue{}, time.Duration(*timeoutSec)*time.Second)
	if err := fut.Wait(); err != nil {
		return errors.Wrap(err, "failed to request: ")
	}
	res, err := fut.Result()
	if err != nil {
		return errors.Wrap(err, "coordinator responds error: ")
	}
	agg, ok := res.(*command.ShowAggregatedValueAck)
	if !ok {
		return fmt.Errorf("invalid CoordinatorStatsAck: %#v", res)
	}

	log.Println(fmt.Sprintf("got %d values", len(agg.AggregatedValues)))
	for name, val := range agg.AggregatedValues {
		log.Printf("[%s] %s\n", name, val)
	}
	return nil
}

func sendShutdown(coordinator *actor.PID) error {
	ctx := actor.EmptyRootContext
	fut := ctx.RequestFuture(coordinator, &command.Shutdown{}, time.Duration(*timeoutSec)*time.Second)
	if err := fut.Wait(); err != nil {
		return errors.Wrap(err, "failed to request: ")
	}
	if _, err := fut.Result(); err != nil {
		return errors.Wrap(err, "coordinator responds error: ")
	}
	log.Println("ok")
	return nil
}

func showStat(coordinator *actor.PID) error {
	var stat command.CoordinatorStatsAck
	if err := requestAsJSON(http.MethodGet, worker.APIPathStats, &command.CoordinatorStats{}, &stat); err != nil {
		return err
	}

	printStat(&stat)

	return nil
}

func watch(coordinator *actor.PID) error {
	for {
		var stat command.CoordinatorStatsAck
		if err := requestAsJSON(http.MethodGet, worker.APIPathStats, &command.CoordinatorStats{}, &stat); err != nil {
			return err
		}

		printStat(&stat)

		if stat.StatsCompleted() {
			log.Println("")
			log.Println("completed")
			break
		}

		time.Sleep(time.Duration(*watchDuration) * time.Millisecond)
	}

	return nil
}

func requestAsJSON(method string, apiPath string, req interface{}, res interface{}) error {
	u, err := url.Parse(fmt.Sprintf("http://%s", *masterHost)) // TODO: https?
	if err != nil {
		return errors.Wrap(err, "invalid master host")
	}
	u.Path = path.Join(apiPath)

	var body io.Reader
	if req != nil {
		b, err := json.Marshal(req)
		if err != nil {
			return errors.Wrap(err, "failed to marshal json")
		}
		body = bytes.NewReader(b)
	}
	request, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return errors.Wrap(err, "failed to new request")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read body")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: status=%d body=%s", response.StatusCode, string(resBody))
	}

	if res != nil {
		if err := json.Unmarshal(resBody, res); err != nil {
			return fmt.Errorf("failed to parse response: %s", resBody)
		}
	}

	return nil
}

func printStat(s *command.CoordinatorStatsAck) {
	sb := strings.Builder{}
	sb.WriteString("state=")
	sb.WriteString(s.State)
	sb.WriteString(" superstep=")
	sb.WriteString(strconv.FormatUint(s.SuperStep, 10))
	sb.WriteString(" active=")
	sb.WriteString(strconv.FormatUint(s.NrOfActiveVertex, 10))
	sb.WriteString(" sent=")
	sb.WriteString(strconv.FormatUint(s.NrOfSentMessages, 10))
	log.Print(sb.String())
}

func load(coordinator *actor.PID, ids []string) error {
	ctx := actor.EmptyRootContext
	wg := &sync.WaitGroup{}

	for _, vid := range ids {
		id := vid
		wg.Add(1)
		go func() {
			defer wg.Done()

			fut := ctx.RequestFuture(coordinator, &command.LoadVertex{
				VertexId: id,
			}, time.Duration(*timeoutSec)*time.Second)

			if err := fut.Wait(); err != nil {
				log.Printf("ERROR! %v\n", err.Error())
				return
			}

			res, err := fut.Result()
			if err != nil {
				log.Printf("ERROR! %v\n", err.Error())
				return
			}

			ack, ok := res.(*command.LoadVertexAck)
			if !ok {
				log.Printf("ERROR! unexpected message %#v\n", res)
				return
			}

			if ack.Error != "" {
				log.Printf("ERROR! %s\n", ack.Error)
				return
			}

			log.Printf("complete: veretex id is %s\n", ack.VertexId)
		}()
	}

	wg.Wait()
	log.Println("done")

	return nil
}

func startSuperStep(coordinator *actor.PID) error {
	actor.EmptyRootContext.Send(coordinator, &command.StartSuperStep{})
	return watch(coordinator)
}
