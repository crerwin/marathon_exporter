FROM scratch
MAINTAINER Chris Erwin "https://github.com/crerwin/marathon_exporter/"

ADD bin/marathon_exporter /marathon_exporter
ENTRYPOINT ["/marathon_exporter"]

EXPOSE 9088
