FROM golang:alpine AS builder

ENV GO111MODULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

ARG WORK=/build

WORKDIR $WORK

ADD . .
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go mod download
RUN go build -o lottery

FROM szyhf/alpine-glibc:latest
WORKDIR /runner
COPY --from=builder /build .
RUN chmod u+x lottery
ENTRYPOINT [ "./lottery" ]
