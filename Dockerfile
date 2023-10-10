FROM golang:1.18-alpine3.16 AS builder

WORKDIR /build

ENV GOPROXY https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o serve ./cmd/server/

FROM alpine:3.16 AS final

WORKDIR /app

COPY --from=builder /build/serve /app/serve

ENTRYPOINT ["/app/serve"]
