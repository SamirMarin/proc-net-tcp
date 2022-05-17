# proc-net-tcp
Reads /proc/net/tcp every 10 seconds and outputs any new connections


## How to build and run 

Uses golang 1.17

Tested on amazon/linux ec2

### run using docker

```
make docker-run
```

### run directly

```
make build

./main
```

### run test

```
make test
```


## Server

- listens on port 8080 

- Reads [/proc/net/tcp](https://github.com/SamirMarin/proc-net-tcp/blob/ebd15e9edc16222ac80be6862845460e519b7165/handlers/tcp.go#L11) every [10 seconds](https://github.com/SamirMarin/proc-net-tcp/blob/ebd15e9edc16222ac80be6862845460e519b7165/handlers/tcp.go#L12)

- Blocks source IP when port scan is [detected](https://github.com/SamirMarin/proc-net-tcp/blob/ebd15e9edc16222ac80be6862845460e519b7165/pkg/tcp/tcp.go#L203)

### endpoints
- /metrics ---> [Counter - number of new connections](https://github.com/SamirMarin/proc-net-tcp/blob/ebd15e9edc16222ac80be6862845460e519b7165/pkg/tcp/tcp.go#L135)

- /health ---> health check
