package worker

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// CommonConfig is set of environments for both of worker and master
type CommonConfig struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"INFO"`
}

// WorkerEnv is set of environments for workers
type WorkerEnv struct {
	CommonConfig
	ListenAddress string `envconfig:"LISTEN_ADDR" default:"127.0.0.1:8801"`
}

// MasterEnv is set of environments for workers
type MasterEnv struct {
	CommonConfig
	ListenAddress   string `envconfig:"LISTEN_ADDR" default:"127.0.0.1:8801"`
	WorkerAddresses string `envconfig:"WORKERS" default:""`
}

// ReadEnv reads configuration from env
func ReadEnv(prefix string) (interface{}, error) {
	role := os.Getenv("ROLE")

	switch role {
	case "worker":
		var workerEnv WorkerEnv
		if err := envconfig.Process(prefix, &workerEnv); err != nil {
			return nil, errors.Wrap(err, "failed to parse env for master: ")
		}
		return &workerEnv, nil

	case "master":
		var masterEnv MasterEnv
		if err := envconfig.Process(prefix, &masterEnv); err != nil {
			return nil, errors.Wrap(err, "failed to parse env for master: ")
		}
		return &masterEnv, nil

	default:
		return nil, fmt.Errorf("invalid ROLE(%s), it must be one of worker, master", role)
	}
}

// SetupLogger configures logger
func (cc *CommonConfig) Logger() *logrus.Logger {
	lv, err := logrus.ParseLevel(cc.LogLevel)
	if err != nil {
		logrus.Warn(err)
	} else {
		logrus.SetLevel(lv)
	}
	return logrus.New()
}
