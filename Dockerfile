FROM golang:1.16.3 as build

WORKDIR /usr/local/go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["uptotg"]
