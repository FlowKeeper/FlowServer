FROM golang:1.17-alpine

COPY . /src
RUN apk add gcc musl-dev
RUN cd /src && go build -o /src/server .

FROM alpine:latest
RUN mkdir /app
COPY --from=0 /src/server /app/server

CMD ["/app/server"]

