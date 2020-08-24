FROM golang

ARG registries

WORKDIR /opt

COPY . .

RUN make --makefile=Makefile.src && \
    make --makefile=Makefile.src install && \
    make --makefile=Makefile.src clean

# empty unless the image is specifically built with it
# the docker plugin install command will set this later if needed
ENV REGISTRIES=${registries}

ENTRYPOINT /usr/libexec/img-authz-plugin
