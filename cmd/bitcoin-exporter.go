package main

import (
    "io/ioutil"
    "log"
    "os"
    "github.com/justinbarrick/bitcoin-exporter"
    "github.com/justinbarrick/bitcoin-exporter/coinbase"
    "github.com/justinbarrick/bitcoin-exporter/metrics"
    "gopkg.in/yaml.v2"
)

func main() {
    config := exporter.Config{}

    config_name := "bitcoin-exporter.yml"
    if len(os.Args) > 1 {
        config_name = os.Args[1]
    }

    data, err := ioutil.ReadFile(config_name)
    if err != nil {
        log.Fatal(err)
    }

    err = yaml.Unmarshal(data, &config)
    if err != nil {
        log.Fatal(err)
    }

    go coinbase.Monitor(5, &config)
    metrics.Init()
}
