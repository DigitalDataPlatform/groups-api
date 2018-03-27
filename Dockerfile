FROM golang:1.10 as builder

WORKDIR /go/src/gitlab.adeo.com/ddp-portal-api
copy . .

RUN go get github.com/apex/log
RUN go get github.com/go-chi/chi
RUN go get github.com/go-chi/render
RUN go get github.com/gobuffalo/envy
RUN go get github.com/dgraph-io/badger
RUN go get github.com/davecgh/go-spew/spew
RUN go get github.com/segmentio/ksuid
RUN go get github.com/pkg/errors

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./api/cmd/server.go

FROM alpine

COPY --from=builder /go/src/gitlab.adeo.com/ddp-portal-api/server /

CMD ["./server"]
