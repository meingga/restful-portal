FROM golang:1.21.1-alpine

WORKDIR /app

COPY . .

RUN go build -o restful-portal

EXPOSE 8383

CMD ./restful-portal
