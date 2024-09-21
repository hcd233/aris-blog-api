# 使用官方的Go镜像作为构建阶段
FROM golang:1.22.5-alpine3.20 AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor

RUN go mod tidy

COPY . .

RUN go build -mod=vendor -o /go/bin/aris-blog-api

FROM alpine:latest

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app

COPY --from=builder /go/bin/aris-blog-api /app/aris-blog-api

EXPOSE 8080

# CMD ["/app/aris-blog-api", "server", "start", "--host", "0.0.0.0", "--port", "8080"]

# docker buildx build --platform linux/amd64 -t aris-blog-api:latest .
# docker run -d -p 8080:8080 --env-file api.env --name aris-blog-api -t aris-blog-api:latest