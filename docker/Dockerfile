FROM golang:alpine AS builder

ENV GO111MODULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

ARG WORK=/build

WORKDIR $WORK

ADD . .
RUN go mod download
RUN go build -o lottery

FROM szyhf/alpine-glibc:latest
WORKDIR /runner
COPY --from=builder /build/lottery .
RUN chmod u+x lottery
ENTRYPOINT [ "./lottery" ]
