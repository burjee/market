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

RUN go build -o bin/server-match .

CMD ["/app/bin/server-match"]

FROM scratch

COPY --from=builder /app/bin/server-match /app/bin/server-match

# ENTRYPOINT
CMD ["/app/bin/server-match"]