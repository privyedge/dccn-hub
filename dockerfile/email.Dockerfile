FROM golang:1.11.4-alpine as builder

RUN apk add -U --no-cache ca-certificates

WORKDIR /go/src/github.com/Ankr-network/dccn-hub

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -i -o cmd/email app-dccn-email/main.go

FROM scratch

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/src/github.com/Ankr-network/dccn-hub/cmd/email /

CMD ["/email"]
