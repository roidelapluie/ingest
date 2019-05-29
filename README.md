# ingest

Backfill prometheus.

## Build

`
go build
`

## args

```
  -block-range int
        block range (default 1048576)
  -input-file string
        json input file (default "input")
  -output-dir string
        directory where to put the TSDB (default "output")
```

## Use

1. Create and populate TSDB

```
(echo "["
for i in $(seq $(date +%s -d "1 month ago")000 1000 $(date +%s)000)
do echo '{"l":{"__name__":"test"}, "t":'$i', "v":'$i'},'
done) | ./ingest -input-file=/dev/stdin -output-dir 1month
```

2. Start prometheus on top of it

```
prometheus --storage.tsdb.retention.time=99y --storage.tsdb.path=1month --config.file=$(mktemp)
```

## Json format

```
[
    {"l":{"__name__":"test"}, "t":1, "v":1}
]
```

Array of metrics. Metric being objects with 3 fields: labels, timestamps (epoch
milliseconds) and value.

There is no need to close the json:

```
[
    {"l":{"__name__":"test"}, "t":1, "v":1},
    {"l":{"__name__":"test"}, "t":2, "v":2},
    {"l":{"__name__":"test"}, "t":3, "v":3},
    {"l":{"__name__":"test"}, "t":4, "v":4},
```
