FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o /server ./cmd/server

FROM alpine:3.21
COPY --from=builder /server /server
ENTRYPOINT ["/server"]
