FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/ao3-scrapper ./cmd/ao3-scrapper/main.go

FROM alpine

COPY --from=builder /go/bin/ao3-scrapper /go/bin/ao3-scrapper

ENTRYPOINT ["/go/bin/ao3-scrapper"]