FROM golang:1.15-alpine AS base
WORKDIR $GOPATH/src/github.com/gaoxinge/website-v2-operator/
COPY . .
ENV GOPROXY=https://goproxy.io,direct
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /cmd/main cmd/main.go

FROM alpine
COPY --from=base /cmd/main /cmd/main
ENTRYPOINT ["/cmd/main"]