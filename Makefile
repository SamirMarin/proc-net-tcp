PROC_NET_TCP_SEVER=tcp
PROC_NET_TCP=proc-net-tcp
TAG=0.0.1

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
		--cap-add NET_ADMIN \
		--net=host \
		${PROC_NET_TCP}:${TAG}
