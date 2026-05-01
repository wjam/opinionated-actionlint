FROM golang:1.26@sha256:595c7847cff97c9a9e76f015083c481d26078f961c9c8dca3923132f51fe12f1 as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go run ./build

FROM koalaman/shellcheck-alpine:v0.11.0@sha256:9955be09ea7f0dbf7ae942ac1f2094355bb30d96fffba0ec09f5432207544002 as shellcheck

FROM alpine:3.23@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11
COPY --from=builder /src/bin/app /usr/local/bin/opinionated-actionlint
COPY --from=shellcheck /bin/shellcheck /usr/local/bin/shellcheck

USER guest
ENTRYPOINT ["/usr/local/bin/opinionated-actionlint"]
