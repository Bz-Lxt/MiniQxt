FROM golang:1.22-alpine AS builder
WORKDIR /src
ENV GOPROXY=https://goproxy.cn,direct \
    GOTOOLCHAIN=auto
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM alpine:3.20
WORKDIR /app
ENV TZ=Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add --no-cache tzdata ca-certificates wget ffmpeg \
    && mkdir -p /app/media \
    && ffmpeg -y -f lavfi -i "color=c=0x12233f:s=1280x720:d=30" \
       -f lavfi -i "sine=frequency=440:duration=30" \
       -shortest -pix_fmt yuv420p /app/media/onboarding.mp4 \
    && apk del ffmpeg
COPY --from=builder /out/server /app/server
EXPOSE 8080
CMD ["/app/server"]
