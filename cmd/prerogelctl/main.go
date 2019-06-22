package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	case args[0] == "watch":
		err = watch(coordinator)
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

func info(msg string) {
	log.Println("INFO: " + msg)
}

func showStat(coordinator *actor.PID) error {
	stat, err := getStat(actor.EmptyRootContext, coordinator)
	if err != nil {
		return err
	}

	printStat(stat)

	return nil
}

func watch(coordinator *actor.PID) error {
	ctx := actor.EmptyRootContext

	for {
		stat, err := getStat(ctx, coordinator)
		if err != nil {
			return err
		}

		printStat(stat)

		if stat.StatsCompleted() {
			log.Println("")
			info("complete")
			break
		}

		time.Sleep(time.Duration(*watchDuration) * time.Millisecond)
	}

	return nil
}

func getStat(ctx actor.SenderContext, coordinator *actor.PID) (*command.CoordinatorStatsAck, error) {
	fut := ctx.RequestFuture(coordinator, &command.CoordinatorStats{}, time.Duration(*timeoutSec)*time.Second)
	if err := fut.Wait(); err != nil {
		return nil, errors.Wrap(err, "failed to ask coordinator: ")
	}
	res, err := fut.Result()
	if err != nil {
		return nil, errors.Wrap(err, "coordinator responds error: ")
	}
	stat, ok := res.(*command.CoordinatorStatsAck)
	if !ok {
		return nil, fmt.Errorf("invalid CoordinatorStatsAck: %#v", res)
	}
	return stat, nil
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
