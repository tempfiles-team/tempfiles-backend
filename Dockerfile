FROM ubuntu:latest

ARG TARGETPLATFORM

WORKDIR /app

COPY /dist/tempfiles-backend_$TARGETPLATFORM/tempfiles-backend /

EXPOSE 5000

CMD ["./tempfiles-backend"]