FROM golang

ARG registries

WORKDIR /opt

COPY . .

ENV REGISTRIES=${registries}

RUN mkdir -p /usr/libexec/ /run/docker/plugins && make && make config REGISTRIES=$REGISTRIES && make install

ENTRYPOINT /usr/libexec/img-authz-plugin
