FROM crpi-aejpgd9d8unswtqu.cn-shenzhen.personal.cr.aliyuncs.com/immigration/golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/server .
COPY database/migrations ./database/migrations
EXPOSE 8080
CMD ["./server"]
