PROC_NET_TCP_SEVER=tcp

clean:
	rm main

test:
	go test -v ./...


build: test
	go build cmd/${PROC_NET_TCP_SEVER}/main.go
