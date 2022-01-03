ARG GO_VERSION=alpine

FROM golang:${GO_VERSION}

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64 \
    GOPROXY=https://goproxy.cn

# 移动到工作目录：/build
WORKDIR /data

RUN go install github.com/cosmtrek/air@latest