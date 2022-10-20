FROM golang:1.17-alpine3.15 as builder
ENV GOPROXY="https://goproxy.cn" \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on

RUN apk add build-base
WORKDIR /
COPY ./src /src
WORKDIR /src
RUN go build -o /build/app .

FROM alpine:3.15

RUN apk --no-cache add -U tzdata ca-certificates libc6-compat libgcc libstdc++ && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > "/etc/timezone"
COPY --from=builder /build/app /usr/bin/app
COPY ./config /config
WORKDIR /
EXPOSE 3487
ENTRYPOINT [ "app" ]
