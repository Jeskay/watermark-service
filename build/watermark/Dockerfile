FROM golang:1.21-alpine3.19 as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GO_ARCH=amd64 go build -o watermark_service ./cmd/watermark/watermark.go

FROM scratch

COPY --from=builder /app/watermark_service ./
COPY --from=builder /app/var ./var

EXPOSE 9091 9092

#CMD ["./watermark_service", "-config", "./config/watermark_config.yaml"]

CMD ["./watermark_service"]