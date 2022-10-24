FROM golang:1.19 AS builder
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY back/ /go/src/app/back/
COPY env/ /go/src/app/env/
WORKDIR /go/src/app/back
RUN CGO_ENABLED=0 go build -o back

FROM gcr.io/distroless/static-debian11 AS runtime
COPY --from=builder /go/src/app/back/back /usr/local/bin/back
ENTRYPOINT ["/usr/local/bin/back"]
