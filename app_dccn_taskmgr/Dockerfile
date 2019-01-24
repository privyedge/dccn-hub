FROM golang:1.11.4-alpine as builder

WORKDIR /go/src/github.com/Ankr-network/dccn-hub/app_dccn_taskmgr
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app_dccn_taskmgr .

FROM scratch

COPY --from=builder /go/src/github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/app_dccn_taskmgr /
CMD ["/app_dccn_taskmgr"]
