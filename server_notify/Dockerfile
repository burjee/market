FROM golang:1.22.6-alpine AS dev

WORKDIR /app

ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GOCACHE /go-build

COPY go.mod go.sum ./
RUN go mod download

CMD ["go", "run", "."]

FROM dev as builder

COPY . .

RUN go build -o bin/server-notify main.go

CMD ["/app/bin/server-notify"]

FROM scratch

COPY --from=builder /app/bin/server-notify /app/bin/server-notify

# ENTRYPOINT
CMD ["/app/bin/server-notify"]