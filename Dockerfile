# Copyright (c) 2024 FajarLaksono. All Rights Reserved.

FROM golang:1.23.6-alpine
LABEL maintainer="fajrlaksono@gmail.com"

ENV HOME /srv

COPY service $HOME/service
COPY contents/swagger-ui $HOME/contents/swagger-ui
COPY migrations $HOME/migrations
COPY version.json $HOME/version.json
COPY LICENSE $HOME/LICENSE

RUN find $HOME -type d -exec 'chmod' '555' '{}' ';' && \
    find $HOME -type f -exec 'chmod' '444' '{}' ';' && \
    find $HOME -type f -exec 'chown' 'root:root' '{}' ';' && \
    chmod 555 $HOME/service

USER nobody

ENTRYPOINT ["/srv/service"]
