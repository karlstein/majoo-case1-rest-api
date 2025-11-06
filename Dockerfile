FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/blog-api ./cmd/http

FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=builder /bin/blog-api /onagatego
COPY config/.env /.env
EXPOSE 3011
ENV PORT=3011
ENTRYPOINT ["/onagatego", "--env-path", ".env"]


