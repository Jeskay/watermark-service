FROM golang:1.21-alpine3.19 as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GO_ARCH=amd64 go build -o auth_service ./cmd/authentication/auth.go

FROM scratch
COPY --from=builder /app/auth_service ./
COPY --from=builder /app/var ./var

EXPOSE 9021 9022

#CMD ["./auth_service", "-config", "./config/auth_config.yaml"]

CMD ["./auth_service"]