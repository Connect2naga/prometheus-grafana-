# Build the manager binary
FROM golang:1.21 

WORKDIR /go/src/prometheus-grafana

# CMD ["/bin/sh", "-ec", "while :; do echo '.'; sleep 5 ; done"]

CMD go run cmd/sample-service/main.go
 