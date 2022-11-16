FROM ubuntu:latest

ARG TARGETARCH

RUN if [ "$TARGETARCH" = "amd64" ]; then \
    export build_arch="amd64_v1"; \
    else \
    export build_arch="arm64"; \
    fi

ARG build_arch

WORKDIR /app

COPY /dist/tempfiles-backend_linux_$build_arch/tempfiles-backend /
# COPY tempfiles-backend /


EXPOSE 5000

CMD ["./tempfiles-backend"]