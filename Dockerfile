FROM golang:1.17-bullseye AS build
COPY . /src

WORKDIR /src
RUN CGO_ENABLED=0 go build -o /src/server .

FROM alpine:3
RUN mkdir /app
COPY --from=build /src/server /app/server

CMD ["/app/server"]

