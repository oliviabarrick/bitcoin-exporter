package coinbase

import (
    "log"
    "net/http"
    "time"

    gdax "github.com/preichenberger/go-gdax"
    ws "github.com/gorilla/websocket"

    "github.com/justinbarrick/bitcoin-exporter"
    "github.com/justinbarrick/bitcoin-exporter/metrics"
)

func CreateClient(config *exporter.Config) {
    log.Println("Creating coinbase client with key: ", config.Coinbase.Key)

    config.Coinbase.Client = gdax.NewClient(config.Coinbase.Secret, config.Coinbase.Key, config.Coinbase.Passphrase)

    config.Coinbase.Client.HttpClient = &http.Client {
        Timeout: 15 * time.Second,
    }
}

func MonitorPrices(config *exporter.Config) {
    products, err := config.Coinbase.Client.GetProducts()
    if err != nil {
        return
    }

    product_maps := map[string]gdax.Product{}

    product_ids := []string{}
    for _, product := range products {
        product_ids = append(product_ids, product.Id)
        product_maps[product.Id] = product
        log.Println("Product: id: ", product.Id, " base currency: ", product.BaseCurrency, " quote currency: ", product.QuoteCurrency)
    }

    var wsDialer ws.Dialer

    wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
    if err != nil {
        println(err.Error())
    }

    subscribe := gdax.Message{
        Type: "subscribe",
        Channels: []gdax.MessageChannel{
            gdax.MessageChannel{
                Name: "ticker",
                ProductIds: product_ids,
            },
        },
    }

    if err := wsConn.WriteJSON(subscribe); err != nil {
        println(err.Error())
    }

    for true {
        message := gdax.Message{}
        if err := wsConn.ReadJSON(&message); err != nil {
            println(err.Error())
            break
        }

        if message.Type != "ticker" {
            continue
        }

        log.Println("Got ticker:", message.ProductId, message.Side, message.Price)
        metrics.RecordPrice("coinbase", message.ProductId, message.Side, message.Price)
    }
}

func FetchBalance(config *exporter.Config) (err error) {
    accounts, err := config.Coinbase.Client.GetAccounts()
    if err != nil {
        return
    }

    for _, a := range accounts {
        log.Println("Balance:", a.Currency, a.Balance)
        metrics.RecordBalance("coinbase", a.Currency, a.Balance)
    }

    return
}

func Monitor(frequency int, config *exporter.Config) {
    CreateClient(config)

    go func() {
        for {
            MonitorPrices(config)
        }
    }()

    for {
        go func() {
            for i := 0; i < 3; i++ {
                err := FetchBalance(config)
                if err == nil {
                    break
                }
                log.Print(err)
            }

        }()

        time.Sleep(time.Duration(frequency) * time.Second)
    }
}
