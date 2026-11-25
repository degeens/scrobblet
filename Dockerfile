FROM golang:1.27.4-alpine@sha256:70bd7ddea8b6cbe5d6943cde279d6b67feb10cc924fb9100d1d8757c1575b280 AS builder
ARG VERSION=dev
ENV CGO_ENABLED=0 \
    GOOS=linux
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN go build -ldflags="-X main.version=${VERSION}" -o /app ./cmd/scrobblet

FROM alpine:3.24.4@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf AS final
RUN mkdir -p /etc/scrobblet \
    && chown 10001:10001 /etc/scrobblet \
    && chmod 0700 /etc/scrobblet
COPY --from=builder --chown=10001:10001 --chmod=0555 /app /bin/app
USER 10001:10001
EXPOSE 7276
CMD ["/bin/app"]
