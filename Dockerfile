FROM scratch:latest

ARG TARGETARCH

WORKDIR /app

COPY dist/go_multiarch_linux_$TARGETARCH/go_multiarch /

EXPOSE 5000

CMD ["/app/tempfiles-backend"]