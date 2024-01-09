FROM golang:1.21

WORKDIR /

COPY ./bin/app /app
COPY ./templates /templates/

CMD ["/app"]