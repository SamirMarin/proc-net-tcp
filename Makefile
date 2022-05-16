PROC_NET_TCP_SEVER=tcp
PROC_NET_TCP=proc-net-tcp
TAG=0.0.1
PROC_NET_TCP_HOST_PATH=/proc/net/tcp
PROC_NET_TCP_COINTAINER_PATH=/proc-net-tcp/tcp
IPTABLES_HOST_PATH=/usr/sbin/iptables
IPTABLES_COINTAINER_PATH=/proc-net-tcp/iptables

clean:
	rm main

test:
	go test -v ./...

build: test
	go build cmd/${PROC_NET_TCP_SEVER}/main.go

docker-build:
	docker build -t ${PROC_NET_TCP}:${TAG} .

docker-run: docker-build
	docker run \
		-p 8080:8080 \
		-v ${PROC_NET_TCP_HOST_PATH}:${PROC_NET_TCP_COINTAINER_PATH} \
		-v ${IPTABLES_HOST_PATH}:${IPTABLES_COINTAINER_PATH} \
		${PROC_NET_TCP}:${TAG}
