FROM golang:1.2@sha256:3f7444391c51a11a039bf0359ee81cc64e663c17d787ad0e637a4de1a3f62a71 as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go run ./build

FROM koalaman/shellcheck-alpine:v0.10.0@sha256:5921d946dac740cbeec2fb1c898747b6105e585130cc7f0602eec9a10f7ddb63 as shellcheck

FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c
COPY --from=builder /src/bin/app /usr/local/bin/opinionated-actionlint
COPY --from=shellcheck /bin/shellcheck /usr/local/bin/shellcheck

USER guest
ENTRYPOINT ["/usr/local/bin/opinionated-actionlint"]
