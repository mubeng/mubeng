FROM golang:1.24-alpine AS build

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

COPY --from=build /app/bin/mubeng /bin/mubeng
ENV HOME /
ENTRYPOINT ["/bin/mubeng"]
