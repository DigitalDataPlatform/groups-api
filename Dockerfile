FROM golang:1.10

WORKDIR /go/src/portal-api
copy . .

RUN go get github.com/apex/log
RUN go get github.com/go-chi/chi
RUN go get github.com/gobuffalo/envy


RUN go install -v ./...

CMD ["portal-api"]
