FROM golang:1.17.6-alpine3.15

RUN mkdir -p /proc-net-tcp
WORKDIR /proc-net-tcp
COPY . /proc-net-tcp

RUN CGO_ENABLED=0 go build -o /bin/proc-net-tcp cmd/tcp/main.go

EXPOSE 8080
ENTRYPOINT ["/bin/proc-net-tcp"]
