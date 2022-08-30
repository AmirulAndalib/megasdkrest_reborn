FROM ghcr.io/amirulandalib/megarestbase:latest AS builder

ARG CPU_ARCH=amd64
ENV HOST_CPU_ARCH=$CPU_ARCH

# MegaSDK
RUN git clone https://github.com/meganz/sdk.git sdk && cd sdk && \
    sh autogen.sh && \
    ./configure  CPPFLAGS='-fpermissive' --disable-silent-rules \
    --disable-examples --disable-shared --enable-static --without-freeimage && \
    make -j$(getconf _NPROCESSORS_ONLN) && \
    make install

# MegaSDK-Go
RUN mkdir -p /usr/local/go/src/ && cd /usr/local/go/src/ && \
    git clone https://github.com/AmirulAndalib/megasdkgo && \
    cd megasdkgo && rm -rf .git && \
    mkdir include && cp -r /go/sdk/include/* include && \
    mkdir .libs && \
    cp /usr/lib/lib*.a .libs/ && \
    cp /usr/lib/lib*.la .libs/ && \
    go tool cgo megasdkgo.go

# MegaSDK-REST
RUN git clone https://github.com/AmirulAndalib/megasdkrest_reborn temp && cd temp && \
    go build -ldflags "-linkmode external -extldflags '-static' -s -W" . && \
    mkdir -p /go/build/ && mv megasdkrpc ../build/megasdkrest-${HOST_CPU_ARCH}

FROM scratch AS megasdkrest

COPY --from=builder /go/build/ /
