FROM alpine:3

ARG TARGETOS
ARG TARGETARCH

ENV EXT_ENV_LOG_LEVEL=info

COPY dist/$TARGETOS/$TARGETARCH/drone_ext_envs          /bin/

COPY dist/$TARGETOS/$TARGETARCH/drone_ext_envs_client   /bin/

ENTRYPOINT ["/bin/drone_ext_envs"]