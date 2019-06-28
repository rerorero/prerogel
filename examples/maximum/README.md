Maximum Value
====

```$sh
# Start prerogel cluster on docker-compose.
$ docker-compose up
```

```$sh
# You can access prerogel cluster with `prerogelctl`.
$ go get -u github.com/rerorero/prerogel/cmd/prerogelctl

# Load the vertices then start calculation.
$ prerogelctl -host 127.0.0.1:9000 load a b c d

$ prerogelctl -host 127.0.0.1:9000 start

# show aggregator values
# maximum value is calculated as an aggregator value
$ prerogelctl -host 127.0.0.1:9000 agg

# Destroy cluster
$ prerogelctl -host 127.0.0.1:9000 shutdown
```
