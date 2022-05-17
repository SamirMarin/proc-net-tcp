# proc-net-tcp
Reads /proc/net/tcp every 10 seconds and outputs any new connections


## How to build and run 

Program was tested using golang 1.17

Ran on amazon/linux ec2

### run directly

```
make build

./main
```

To run test

```
make test
```


### run using docker

```
make docker-run
```


### Server

the server listens on port 8080 

/metrics ---> Counter - number of new connections

/health ---> health check

server reads `proc/net/tcp` every 10 seconds outputs new connections, blocks IPs where port scan is detected
