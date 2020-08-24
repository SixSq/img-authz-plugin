FROM golang

ARG registries

WORKDIR /opt

COPY . .

RUN make --makefile=Makefile.src && \
    make install && \
    make clean

# empty unless the image is specifically built with it
# the docker plugin install command will set this later if needed
ENV REGISTRIES=${registries}

ENTRYPOINT /usr/libexec/img-authz-plugin
