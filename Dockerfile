FROM golang:latest

RUN mkdir -p /app

WORKDIR /app

ADD . /app

ENV GOOGLE_APPLICATION_CREDENTIALS="./credentials.json"

RUN go build ./main.go

CMD ["./main"]