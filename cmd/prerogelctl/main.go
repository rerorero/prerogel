package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pkg/errors"
	"github.com/rerorero/prerogel/command"
	"github.com/rerorero/prerogel/worker"
)

func main() {
	os.Exit(realMain())
}

var (
	masterHost    = flag.String("host", "", "host address and port of master worker")
	timeoutSec    = flag.Int("timeout", 30, "timeout (seconds)")
	watchDuration = flag.Int("duration", 300, "ping duration (milliseconds)")
)

func realMain() int {
	flag.Parse()
	args := flag.Args()

	// parse flags
	if *masterHost == "" {
		exitErr(errors.New("no host specified"))
	}

	// common setup
	coordinator := &actor.PID{
		Address: *masterHost,
		Id:      worker.CoordinatorActorID,
	}

	// run command
	var err error
	switch {
	case len(args) == 0:
		err = errors.New("no command specified")
	case args[0] == "watch":
		err = watch(coordinator)
	default:
		err = fmt.Errorf("%s - no such command", args[0])
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

func watch(coordinator *actor.PID) error {
	ctx := actor.EmptyRootContext
	cmd := &command.CoordinatorStats{}

	print := func(s *command.CoordinatorStatsAck) {
		sb := strings.Builder{}
		sb.WriteString("superstep=")
		sb.WriteString(strconv.FormatUint(s.SuperStep, 10))
		sb.WriteString(" active=")
		sb.WriteString(strconv.FormatUint(s.NrOfActiveVertex, 10))
		sb.WriteString(" sent=")
		sb.WriteString(strconv.FormatUint(s.NrOfSentMessages, 10))
		log.Print(sb.String())
	}

	for {
		fut := ctx.RequestFuture(coordinator, cmd, time.Duration(*timeoutSec)*time.Second)
		if err := fut.Wait(); err != nil {
			return errors.Wrap(err, "failed to ask coordinator: ")
		}
		res, err := fut.Result()
		if err != nil {
			return errors.Wrap(err, "coordinator responds error: ")
		}
		stat, ok := res.(*command.CoordinatorStatsAck)
		if !ok {
			return fmt.Errorf("invalid CoordinatorStatsAck: %#v", res)
		}

		print(stat)

		if command.StatsCompleted(stat) {
			log.Println("")
			info("complete")
			break
		}

		time.Sleep(time.Duration(*watchDuration) * time.Millisecond)
	}

	return nil
}
