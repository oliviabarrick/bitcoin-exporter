The bitcoin-exporter scrapes the GDAX API for the trading prices of various
cryptocurrencies and your account balances and produces Prometheus metrics that
you can alert on.

# Configuration

Create a `config/bitcoin-exporter.yml` file containing your coinbase GDAX credentials:

```
coinbase:
  key: baoenuth32401923hsnaotehu9123
  secret: 245092394023euenchu/01923==
  passphrase: hello
```

Next, create an alertmanager configuration called `config/alertmanager.yml`:

```
route:
  receiver: 'slack'

receivers:
  - name: 'slack'
    slack_configs:
    - send_resolved: true
      username: 'Bitcoin Bot'
      channel: '#bitcoin'
      api_url: 'your api url'
      text: "{{ range .Alerts }}{{ .Annotations.description }}\n{{ end }}"
```

# Alert writing

To write an alert, add an alert into `config/alerts/`.

# Starting the stack

To start a stack, complete with the bitcoin-exporter, Grafana, Prometheus, and AlertManager, run:

```
make build
make up
```

If you run `git pull` or make any changes, use `make build` to rebuild before `make up`.

Run `make status` to print out the IPs of Grafana and Prometheus.

Run `make logs` to print the logs of the bitcoin-exporter service.

If you receive no balance metrics, make sure that your clock is correct.

# Building a Docker image

If you already run Prometheus, Grafana, etc, then you can just build the docker image:

```
make build-image
```

This will build a Docker image tagged as `bitcoin-exporter` that you can use.

# Building a binary

If you just want the binary to run, you can just build that:

```
make build-bin
```
