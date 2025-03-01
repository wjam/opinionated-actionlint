FROM golang:1.24 as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
ENV CGO_ENABLED 0
RUN go build -o opinionated-actionlint -trimpath -ldflags "-s -w" .

FROM koalaman/shellcheck-alpine:v0.10.0 as shellcheck

FROM alpine:3.21
COPY --from=builder /src/opinionated-actionlint /usr/local/bin/
COPY --from=shellcheck /bin/shellcheck /usr/local/bin/shellcheck

USER guest
ENTRYPOINT ["/usr/local/bin/opinionated-actionlint"]
