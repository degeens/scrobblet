FROM golang:1.25-alpine AS builder
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app ./cmd/api

FROM alpine:3.23 AS final
COPY --from=builder /app /bin/app
EXPOSE 7276
CMD ["bin/app"]