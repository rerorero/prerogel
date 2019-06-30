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

	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/worker"
)

func main() {
	os.Exit(realMain())
}

var (
	masterHost    = flag.String("host", "", "host address and port of master worker")
	watchDuration = flag.Int("duration", 300, "ping duration (milliseconds)")
	degub         = flag.Bool("debug", false, "debug mode")
	// TODO: maxStep         = flag.Int("max-step", 0, "maximum number of supper step, 0 means no limit")
)

func realMain() int {
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()

	// parse flags
	if *masterHost == "" {
		flag.PrintDefaults()
		exitErr(errors.New("no host specified"))
	}

	// run command
	var err error
	switch {
	case len(args) == 0:
		err = errors.New("no command specified")
	case args[0] == "state":
		err = showStat()
	case args[0] == "load":
		if len(args[1:]) > 0 {
			err = loadVertex(args[1:])
		} else {
			err = loadPartition()
		}
	case args[0] == "start":
		err = startSuperStep()
	case args[0] == "watch":
		err = watch()
	case args[0] == "agg":
		err = showAggregatedValue()
	case args[0] == "shutdown":
		err = sendShutdown()
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

func showAggregatedValue() error {
	var agg command.ShowAggregatedValueAck
	if err := requestAsJSON(http.MethodGet, worker.APIPathShowAggregatedValue, nil, &agg); err != nil {
		return err
	}

	log.Println(fmt.Sprintf("got %d values", len(agg.AggregatedValues)))
	for name, val := range agg.AggregatedValues {
		log.Printf("[%s] %s\n", name, val)
	}
	return nil
}

func sendShutdown() error {
	if err := requestAsJSON(http.MethodPost, worker.APIPathShutdown, nil, nil); err != nil {
		return err
	}
	log.Println("ok")
	return nil
}

func showStat() error {
	var stat command.CoordinatorStatsAck
	if err := requestAsJSON(http.MethodGet, worker.APIPathStats, nil, &stat); err != nil {
		return err
	}
	printStat(&stat)
	return nil
}

func watch() error {
	for {
		var stat command.CoordinatorStatsAck
		if err := requestAsJSON(http.MethodGet, worker.APIPathStats, nil, &stat); err != nil {
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

func loadVertex(ids []string) error {
	wg := &sync.WaitGroup{}

	for _, vid := range ids {
		id := vid
		wg.Add(1)
		go func() {
			defer wg.Done()
			var ack command.LoadVertexAck
			if err := requestAsJSON(http.MethodPost, worker.APIPathLoadVertex, &command.LoadVertex{
				VertexId: id,
			}, &ack); err != nil {
				log.Printf("ERROR! %v\n", err.Error())
				return
			}
			log.Printf("complete: veretex id is %s\n", ack.VertexId)
		}()
	}

	wg.Wait()
	log.Println("done")

	return nil
}

func loadPartition() error {
	if err := requestAsJSON(http.MethodGet, worker.APIPathLoadPartitionVertices, nil, nil); err != nil {
		return err
	}

	for {
		var stat command.CoordinatorStatsAck
		if err := requestAsJSON(http.MethodGet, worker.APIPathStats, nil, &stat); err != nil {
			return err
		}

		if stat.State != worker.CoordinatorStateLoadingVertices {
			if stat.State != worker.CoordinatorStateIdle {
				return fmt.Errorf("current coordinator state: %s", stat.State)
			}
			log.Println("finished")
			break
		}

		time.Sleep(time.Duration(*watchDuration) * time.Millisecond)
	}

	return nil

}

func startSuperStep() error {
	if err := requestAsJSON(http.MethodPost, worker.APIPathStartSuperStep, nil, nil); err != nil {
		return err
	}
	return watch()
}
