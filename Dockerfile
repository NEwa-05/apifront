FROM golang:1.19.3 AS builder
WORKDIR /app
COPY apifront.go .
RUN go mod init apifront
RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o apifront .

FROM alpine
WORKDIR /
COPY --from=builder /app/apifront .
COPY front_template/test.html .
CMD ["./apifront"]
