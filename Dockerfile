FROM golang:1.23-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
ARG version

WORKDIR /src
COPY . /src

RUN go build -a -tags netgo -ldflags="-s -w -X main.version=$version" -o dist/$TARGETOS/$TARGETARCH/drone-env-provider ./cmd/server
RUN go build -a -tags netgo -ldflags="-s -w" -o dist/$TARGETOS/$TARGETARCH/drone-env-provider-client ./cmd/client


#FROM debian:bullseye-slim AS with-build
FROM alpine:3


ARG TARGETOS
ARG TARGETARCH

ENV EXT_ENV_LOG_LEVEL=info

COPY --from=builder /src/dist/$TARGETOS/$TARGETARCH/drone-env-provider          /bin/

COPY --from=builder /src/dist/$TARGETOS/$TARGETARCH/drone-env-provider-client   /bin/

ENTRYPOINT ["/bin/drone-env-provider"]


#FROM debian:bullseye-slim
FROM alpine:3

ARG TARGETOS
ARG TARGETARCH

ENV EXT_ENV_LOG_LEVEL=info

COPY dist/$TARGETOS/$TARGETARCH/drone-env-provider          /bin/

COPY dist/$TARGETOS/$TARGETARCH/drone-env-provider-client   /bin/

ENTRYPOINT ["/bin/drone-env-provider"]

