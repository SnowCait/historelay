FROM golang:1.22

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY main.go ./
RUN go build -v -o /usr/local/bin/app ./...

CMD [ "app" ]
