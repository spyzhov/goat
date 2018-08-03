# About:Goat

Goat is `golang: application template`

Generate template of your application, via chosen templates:
* RabbitMQ - consumer or publisher via [https://github.com/streadway/amqp](https://github.com/streadway/amqp);
* Postgres - database via [https://github.com/go-pg](https://github.com/go-pg);
* Echo - http server via [https://github.com/labstack/echo](https://github.com/labstack/echo);
* Prometheus - prometheus endpoint `GET /metrics` via [https://github.com/prometheus/client_golang](https://github.com/prometheus/client_golang);
* Babex - babex-node for pipeline [https://github.com/matroskin13/babex](https://github.com/matroskin13/babex);

# Installation

Just get it

```bash
go get -u github.com/spyzhov/goat
```

# Basic usage

Change dir to targeted, run `goat` and answer for several questions.

```
user@user:/go/src/github.com/spyzhov/example$ goat
Project path [/go/src/github.com/spyzhov/example]? [Y/n]: y 
Project name [example]? [Y/n]: y
Repository name [github.com/spyzhov/example]? [Y/n]: 
Use Postgres connection (github.com/go-pg)? [y/N]: y 
Use HTTP server (github.com/labstack/echo)? [y/N]: y
Use Prometheus (github.com/prometheus/client_golang)? [y/N]: 
Use Babex-service (github.com/matroskin13/babex)? [y/N]: 
Use RMQ-consumers (github.com/streadway/amqp)? [y/N]: 
Use RMQ-publishers (github.com/streadway/amqp)? [y/N]:
```

# Use [dep](https://github.com/golang/dep)

The best and quickest way to start:
```
user@user:/go/src/github.com/spyzhov/example$ goat
...
user@user:/go/src/github.com/spyzhov/example$ dep init
...
... Profit!
```

# License

MIT licensed. See the [LICENSE](LICENSE) file for details.
