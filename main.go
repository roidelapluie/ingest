package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/tsdb"
	"github.com/prometheus/tsdb/labels"
)

var (
	dir        = flag.String("output-dir", "output", "directory where to put the TSDB")
	input      = flag.String("input-file", "input", "json input file")
	blockRange = flag.Int64("block-range", 1024*1024, "block range")
)

type Metric struct {
	Labels    map[string]string `json:"l"`
	Timestamp int64             `json:"t"`
	Value     float64           `json:"v"`
}

func main() {
	flag.Parse()
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

	opts := &tsdb.Options{
		BlockRanges:    []int64{*blockRange},
		WALSegmentSize: 0,
		NoLockfile:     true,
	}
	db, err := tsdb.Open(*dir, logger, prometheus.DefaultRegisterer, opts)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.EnableCompactions()
	app := db.Appender()
	defer app.Commit()

	f, err := os.Open(*input)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	dec := json.NewDecoder(r)
	_, err = dec.Token()
	if err != nil {
		panic(err)
	}
	for dec.More() {
		var m Metric
		err = dec.Decode(&m)
		if err != nil {
			logger.Log(err)
			continue
		}
		fmt.Printf("%v %v %v\n", m.Labels, m.Timestamp, m.Value)
		_, err := app.Add(labels.FromMap(m.Labels), m.Timestamp, m.Value)
		if err != nil {
			panic(err)
		}
	}
	_, err = dec.Token()
	if err != nil {
		logger.Log(err)
	}
}
