FROM ubuntu:latest

ARG TARGETARCH

WORKDIR /app

COPY /dist/tempfiles-backend_linux_arm_7/tempfiles-backend /

EXPOSE 5000

CMD ["./tempfiles-backend"]