
FROM golang:latest

LABEL "name"="go actions"
LABEL "version"="1.0.0"
LABEL "maintainer"="Michel Courtine <github@michak.net>"

LABEL "com.github.actions.icon"="code"
LABEL "com.github.actions.color"="violet"
LABEL "com.github.actions.name"="Go"
LABEL "com.github.actions.description"="Runs go format, golint or go test. Possible args: format, test, build"

ENV DEP_VERSION="v0.5.0"

RUN curl -fL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/${DEP_VERSION}/dep-linux-amd64 \
 && chmod +x /usr/local/bin/dep

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
