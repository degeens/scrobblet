FROM golang:1.27.1-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS builder
ARG VERSION=dev
ENV CGO_ENABLED=0 \
    GOOS=linux
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN go build -ldflags="-X main.version=${VERSION}" -o /app ./cmd/scrobblet

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b AS final
RUN mkdir -p /etc/scrobblet \
    && chown 10001:10001 /etc/scrobblet \
    && chmod 0700 /etc/scrobblet
COPY --from=builder --chown=10001:10001 --chmod=0555 /app /bin/app
USER 10001:10001
EXPOSE 7276
CMD ["/bin/app"]
