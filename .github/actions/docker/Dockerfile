FROM docker:stable

LABEL "name"="Docker actions"
LABEL "version"="1.0.0"
LABEL "maintainer"="Michel Courtine <github@michak.net>"

LABEL "com.github.actions.icon"="package"
LABEL "com.github.actions.color"="violet"
LABEL "com.github.actions.name"="Docker Publish"
LABEL "com.github.actions.description"="docker tag & push, after a Docker Login"

RUN apk add --no-cache ca-certificates bash git curl \
  && rm -rf /var/cache/apk/*

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
