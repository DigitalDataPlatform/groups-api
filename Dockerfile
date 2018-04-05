FROM golang:1.10 as builder

WORKDIR /go/src/gitlab.adeo.com/ddp-portal-api
COPY . .

# Go dep!
RUN go get -u github.com/golang/dep/...
RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./api/cmd/server.go

FROM alpine

COPY --from=builder /go/src/gitlab.adeo.com/ddp-portal-api/server /

CMD ["./server"]
