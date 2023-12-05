FROM golang:1.21.0 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GO_ARCH=amd64 go build -o picture_service ./cmd/picture/picture.go

FROM scratch
COPY --from=builder /app/picture_service ./
COPY --from=builder /app/var ./var

EXPOSE 8081 8082

#CMD ["./picture_service", "-config", "./config/picture_config.yaml"]

CMD ["./picture_service"]