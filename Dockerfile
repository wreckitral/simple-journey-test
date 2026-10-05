# builder stage
FROM golang:1.27 AS builder

WORKDIR /app/

COPY go.mod main.go ./

ARG VERSION=dev

ENV CGO_ENABLED=0

RUN go build -ldflags "-X main.version=${VERSION} -s -w" -o /out/app .

# run stage
FROM scratch

WORKDIR /app

COPY --from=builder /out/app /app/app

EXPOSE 8080

ENTRYPOINT ["/app/app"]
