package exporter

import (
    gdax "github.com/preichenberger/go-gdax"
)

type Config struct {
    Coinbase CoinbaseConfig
}

type CoinbaseConfig struct {
    Key string `yaml:"key"`
    Secret string `yaml:"secret"`
    Passphrase string `yaml:"passphrase"`
    Client *gdax.Client
}
