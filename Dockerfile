FROM golang:1.25@sha256:6bac879c5b77e0fc9c556a5ed8920e89dab1709bd510a854903509c828f67f96 as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go run ./build

FROM koalaman/shellcheck-alpine:v0.11.0@sha256:9955be09ea7f0dbf7ae942ac1f2094355bb30d96fffba0ec09f5432207544002 as shellcheck

FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1
COPY --from=builder /src/bin/app /usr/local/bin/opinionated-actionlint
COPY --from=shellcheck /bin/shellcheck /usr/local/bin/shellcheck

USER guest
ENTRYPOINT ["/usr/local/bin/opinionated-actionlint"]
