FROM golang:1.18.3-alpine3.16
WORKDIR /app
EXPOSE 5000

COPY . .

RUN go mod download
RUN go build -o ./app cmd/vpiska/main.go

CMD [ "./app" ]