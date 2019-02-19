# About:Goat

Goat is a simple `golang: application template` builder.

Generate template of your application, via chosen templates:
* RabbitMQ - consumer or publisher via [https://github.com/streadway/amqp](https://github.com/streadway/amqp);
* Postgres - database via [https://github.com/go-pg](https://github.com/go-pg);
* MySQL - database via [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql);
* HTTP - http server via [https://golang.org/pkg/net/http](https://golang.org/pkg/net/http/);
* Prometheus - prometheus endpoint `GET /metrics` via [https://github.com/prometheus/client_golang](https://github.com/prometheus/client_golang);

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
Repository name [github.com/spyzhov/example]? [Y/n]: y
Use Postgres connection (github.com/go-pg)? [y/N]: y
Use Postgres migrations (github.com/go-pg/migrations)? [y/N]: y
Use MySQL connection (github.com/go-sql-driver/mysql)? [y/N]: y
Use MySQL migrations (github.com/rubenv/sql-migrate)? [y/N]: y
Use HTTP server (net/http)? [y/N]: y
Use Prometheus (github.com/prometheus/client_golang)? [y/N]: y
Use RMQ-consumer (github.com/streadway/amqp)? [y/N]: y
Use RMQ-publisher (github.com/streadway/amqp)? [y/N]: y
```

And you will got :

```
├── app
│   ├── app.go
│   ├── config.go
│   ├── consumer.go
│   ├── http.go
│   ├── logger.go
│   └── publish.go
├── Dockerfile
├── Gopkg.toml
├── main.go
├── migrations
│   ├── 01_init.go
│   ├── mysql
│   │   └── 1-init.sql
│   ├── mysql.go
│   └── postgres.go
├── README.md
└── signals
    └── signals.go
```

# Use [dep](https://github.com/golang/dep)

The best and quickest way to start:
```
user@user:/go/src/github.com/spyzhov/example$ goat
...
user@user:/go/src/github.com/spyzhov/example$ dep ensure
...
... Profit!
```

# License

MIT licensed. See the [LICENSE](LICENSE) file for details.

# TODO

- [ ] validation for path;
- [x] add versions for libs;
- [ ] add choice for http clients (native/fasthttp/echo/etc.);
- [x] remove code "noodles" - make module append methods;
- [ ] normalize migrations;
- [ ] add Redis support `"github.com/gomodule/redigo/redis"`;
- [x] add context & WaitGroups;
