# 阶段一
FROM golang:1.20.6-alpine as builder
WORKDIR /build
ENV GOPROXY=https://goproxy.cn,direct \
    GOPRIVATE="" \
    GONOSUMDB="" \
    GOSUMDB="sum.golang.google.cn" \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
COPY . .
RUN go mod tidy && go build -o app .

# 阶段二
FROM debian:stable-slim
WORKDIR /app
ENV NACOS_USE=false \
    NACOS_NAMESPACE="_NACOS_NAMESPACE" \
    NACOS_SERVER="_NACOS_SERVER" \
    NACOS_USERNAME="_NACOS_USERNAME" \
    NACOS_PASSWORD="_NACOS_PASSWORD"
COPY --from=builder /build/app .
COPY --from=builder /build/log log
COPY --from=builder /build/nacos.yaml .
ENTRYPOINT ["/app/app"]
