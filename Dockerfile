FROM golang:1.15.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o goapp .

################################################

FROM alpine:3.12

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime && \
    echo "Asia/Bangkok" > /etc/timezone && \
    apk del tzdata

WORKDIR /app

COPY --from=builder /app/goapp .

CMD ["/app/goapp"]  
