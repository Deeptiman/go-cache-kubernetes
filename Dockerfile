FROM golang:latest

RUN mkdir -p $GOPATH/src/github.com/go-cache-kubernetes

WORKDIR $GOPATH/src/github.com/go-cache-kubernetes

COPY . $GOPATH/src/github.com/go-cache-kubernetes

RUN go build -a -installsuffix cgo -o go-cache-kubernetes .

CMD ["chmod +x go-cache-kubernetes"]

CMD ["./go-cache-kubernetes"]

EXPOSE 5000