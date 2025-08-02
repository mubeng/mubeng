FROM golang:1.23-alpine AS build

ARG VERSION

LABEL description="An incredibly fast proxy checker & IP rotator with ease."
LABEL repository="https://github.com/mubeng/mubeng"
LABEL maintainer="dwisiswant0"

WORKDIR /app
COPY ./go.mod .
RUN go mod download

COPY . .
RUN go build -ldflags "-s -w -X github.com/mubeng/mubeng/common.Version=${VERSION}" \
	-o ./bin/mubeng .

FROM alpine:latest

COPY --from=build /app/bin/mubeng /usr/local/bin/mubeng

RUN addgroup -g "2000" mubeng && \
	adduser -g "mubeng" -G "mubeng" -u "1000" -h "/app" -D mubeng

USER mubeng:mubeng

ENTRYPOINT ["mubeng"]
