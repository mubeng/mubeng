FROM golang:1.14.2-alpine3.11 as build

LABEL description="An HTTP proxy server as IP rotation, and proxy checker, all applicable HTTP/S and SOCKS5 protocols with ease."
LABEL repository="https://github.com/kitabisa/mubeng"
LABEL maintainer="dwisiswant0"

WORKDIR /app
COPY ./go.mod .
RUN go mod download

COPY . .
RUN go build -o ./bin/mubeng ./cmd/mubeng 

FROM alpine:latest

COPY --from=build /app/bin/mubeng /bin/mubeng
ENV HOME /
ENTRYPOINT ["/bin/mubeng"]