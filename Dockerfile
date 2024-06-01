FROM golang:1.22-alpine as builder

WORKDIR /app
COPY . ./

RUN go build -o ./auth-api ./cmd/

FROM alpine as resolver
COPY --from=builder /app/auth-api /bin/auth-api
ENTRYPOINT ["/bin/auth-api", "-c", "/etc/cs/auth-config.json"]
