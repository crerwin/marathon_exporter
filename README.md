# Marathon Prometheus Exporter

A [Prometheus](http://prometheus.io) metrics exporter for the [Marathon](https://mesosphere.github.io/marathon) Mesos framework.

This exporter exposes Marathon's Codahale/Dropwizard metrics via its `/metrics` endpoint. To learn more, visit the [Marathon metrics doc](http://mesosphere.github.io/marathon/docs/metrics.html).

This was forked from the unmaintained github.com/gettyimages/marathon_exporter and modified to support scraping Marathon metrics through the DC/OS admin router.

## Getting

```sh
$ go get github.com/crerwin/marathon_exporter
```

*\-or-*

```sh
$ docker pull crerwin/marathon_exporter
```

## Using

```sh
Usage of marathon_exporter:
  -marathon.uri string
        URI of Marathon (default "http://marathon.mesos:8080")
  -dcos.token string
        Bearer token for authenticating to the DCOS API
  -web.listen-address string
        Address to listen on for web interface and telemetry. (default ":9088")
  -web.telemetry-path string
        Path under which to expose metrics. (default "/metrics")
  -log.format value
        If set use a syslog logger or JSON logging. Example: logger:syslog?appname=bob&local=7 or logger:stdout?json=true. Defaults to stderr.
  -log.level value
        Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]. (default info)
```

### Docker run example
```
docker run -p 9088:9088 crerwin/marathon_exporter:0.3.0 --marathon.uri=https://mydcoscluster.company.com/marathon/ --dcos.token="$(dcos config show core.dcos_acs_token)"
```
