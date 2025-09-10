FROM golang:1.25.1-alpine3.22 AS build
RUN apk add --no-cache git ca-certificates
RUN go install github.com/pressly/goose/v3/cmd/goose@latest


FROM alpine:3.22
RUN apk add --no-cache ca-certificates libpq
COPY --from=build /go/bin/goose /usr/local/bin/goose
ENTRYPOINT ["/usr/local/bin/goose"]
CMD ["--help"]