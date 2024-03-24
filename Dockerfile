FROM golang:1.22-alpine as golang

RUN apk update && apk add --no-cache git make gcc musl-dev

WORKDIR /build
COPY / .

RUN go mod download
RUN go mod verify

RUN make build

FROM busybox:1.35.0-uclibc as busybox
FROM gcr.io/distroless/base

COPY --from=busybox /bin/sh /bin/sh
COPY --from=busybox /bin/mkdir /bin/mkdir
COPY --from=busybox /bin/cat /bin/cat
COPY --from=busybox /bin/ls /bin/ls

WORKDIR /build
COPY --from=golang /build/bin/image-processor bin/image-processor
RUN mkdir -p public/images
# RUN echo ${PWD} && ls -lR

EXPOSE 8080
ENTRYPOINT ["/build/bin/image-processor"]