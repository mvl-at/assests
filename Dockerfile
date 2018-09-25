FROM golang:1.11-alpine3.7 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
# ENV GO111MODULE=on

WORKDIR /go/src/github.com/mvl-at/assets
RUN apk add --no-cache \
    git \
    musl-dev
COPY . /go/src/github.com/mvl-at/assets
RUN go get ./...
RUN go install -ldflags '-s -w' ./cmd/assserve

# ---

FROM scratch
COPY --from=build /go/bin/assserve /assets
WORKDIR /assets-data
VOLUME  /assets-data
EXPOSE  7302
ENTRYPOINT [ "/assets" ]
