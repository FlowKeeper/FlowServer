FROM golang:1.17-bullseye

COPY . /src
RUN cd /src && CGO_ENABLED=0 go build -o /src/server .

FROM alpine:latest
RUN mkdir /app
COPY --from=0 /src/server /app/server

CMD ["/app/server"]

