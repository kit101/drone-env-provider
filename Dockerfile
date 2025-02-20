FROM golang:1.23-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
ARG version

WORKDIR /src
COPY . /src

RUN go build -a -tags netgo -ldflags="-s -w -X main.version=$version" -o dist/$TARGETOS/$TARGETARCH/drone_ext_envs ./cmd/server
RUN go build -a -tags netgo -ldflags="-s -w" -o dist/$TARGETOS/$TARGETARCH/drone_ext_envs_client ./cmd/client


#FROM debian:bullseye-slim AS with-build
FROM alpine:3


ARG TARGETOS
ARG TARGETARCH

ENV EXT_ENV_LOG_LEVEL=info

COPY --from=builder /src/dist/$TARGETOS/$TARGETARCH/drone_ext_envs          /bin/

COPY --from=builder /src/dist/$TARGETOS/$TARGETARCH/drone_ext_envs_client   /bin/

ENTRYPOINT ["/bin/drone_ext_envs"]


#FROM debian:bullseye-slim
FROM alpine:3

ARG TARGETOS
ARG TARGETARCH

ENV EXT_ENV_LOG_LEVEL=info

COPY dist/$TARGETOS/$TARGETARCH/drone_ext_envs          /bin/

COPY dist/$TARGETOS/$TARGETARCH/drone_ext_envs_client   /bin/

ENTRYPOINT ["/bin/drone_ext_envs"]

