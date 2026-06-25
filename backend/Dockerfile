FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /md-note-api ./cmd/api

FROM alpine:3.20
COPY --from=builder /md-note-api /md-note-api
EXPOSE 8080
ENTRYPOINT ["/md-note-api"]
