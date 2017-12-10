package metrics

import (
    "log"
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type BitcoinPrice struct {
    SellPrice float64
    SellPriceSubtotal float64

    BuyPrice float64
    BuyPriceSubtotal float64

    CurrentBalance float64
    CurrentBalanceUsd float64

    Exchange string
    Currency string
}

var (
    ticker_metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Namespace: "bitcoin_exporter",
        Name: "ticker",
        Help: "The current price of currencies.",
    }, []string{"exchange", "currency", "side"}, )

    current_balance_metric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Namespace: "bitcoin_exporter",
        Name: "current_balance",
        Help: "The current wallet balance in the original currency.",
    }, []string{"exchange", "currency"}, )
)

func Init() {
    prometheus.MustRegister(ticker_metric)
    prometheus.MustRegister(current_balance_metric)

    http.Handle("/metrics", promhttp.Handler())
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func RecordPrice(exchange string, currency string, side string, price float64) {
    ticker_metric.WithLabelValues(exchange, currency, side).Set(price)
}

func RecordBalance(exchange string, currency string, balance float64) {
    current_balance_metric.WithLabelValues(exchange, currency).Set(balance)
}
