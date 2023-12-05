FROM golang:1.21.0 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GO_ARCH=amd64 go build -o watermark_service ./cmd/watermark/watermark.go

FROM scratch
COPY --from=builder /app/watermark_service ./

EXPOSE 9091 9092

#CMD ["./watermark_service", "-config", "./config/watermark_config.yaml"]

CMD ["./watermark_service"]