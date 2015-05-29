FROM busybox:latest
MAINTAINER etworker

ADD idxgen_server_linux_64bit /idxgen_server

EXPOSE 5678

CMD ["/idxgen_server"]
