SSSP (Single Source Shortest Path)
====

## Build a Prerogel Cluster
#### Build using Docker
```$sh
# Start prerogel cluster on docker-compose.
$ docker-compose up
```

#### Build Locally
run workers (in following example, two workers run)
```$sh
$ export GO111MODULE=on
$ ROLE=worker LISTEN_ADDR=127.0.0.1:8801 go run main.go plugin.go
```
```$sh
$ export GO111MODULE=on
$ ROLE=worker LISTEN_ADDR=127.0.0.1:8802 go run main.go plugin.go
```
start master worker
```$sh
$ export GO111MODULE=on
$ ROLE=master LISTEN_ADDR=127.0.0.1:8803 WORKERS=127.0.0.1:8801,127.0.0.1:8802 PARTITIONS=4 API_PORT=9000 go run main.go plugin.go
```


## Run Computation

```$sh
# You can access prerogel cluster using `prerogelctl` command.
$ go get -u github.com/rerorero/prerogel/cmd/prerogelctl

# Load the vertices then start calculation.
$ prerogelctl -host 127.0.0.1:9000 load

$ prerogelctl -host 127.0.0.1:9000 start

# TODO: Get vertex value somehow

# Destroy cluster
$ prerogelctl -host 127.0.0.1:9000 shutdown
```
