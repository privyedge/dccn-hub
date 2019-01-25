FROM golang:1.11.4-alpine as builder

WORKDIR /go/src/github.com/Ankr-network/dccn-hub/api
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i  -o api ./main.go

FROM scratch

COPY --from=builder /go/src/github.com/Ankr-network/dccn-hub/api/api /
CMD ["/api"]
