FROM golang:1.27.2-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS builder
ARG VERSION=dev
ENV CGO_ENABLED=0 \
    GOOS=linux
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN go build -ldflags="-X main.version=${VERSION}" -o /app ./cmd/scrobblet

FROM alpine:3.24.3@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a AS final
RUN mkdir -p /etc/scrobblet \
    && chown 10001:10001 /etc/scrobblet \
    && chmod 0700 /etc/scrobblet
COPY --from=builder --chown=10001:10001 --chmod=0555 /app /bin/app
USER 10001:10001
EXPOSE 7276
CMD ["/bin/app"]
