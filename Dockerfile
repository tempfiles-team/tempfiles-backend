FROM ubuntu:latest

WORKDIR /app

COPY tempfiles-backend /

EXPOSE 5000

CMD ["./tempfiles-backend"]