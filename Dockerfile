FROM docker.io/library/alpine:3.23 as runtime

RUN \
  apk add --update --no-cache \
  bash \
  curl \
  ca-certificates \
  tzdata

ENTRYPOINT ["provider-s3"]
CMD ["operator"]
COPY provider-s3 /usr/bin/

USER 65536:0
