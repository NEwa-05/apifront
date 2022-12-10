FROM golang:1.19.3 AS builder
WORKDIR /app
COPY postdata_backend/postdata_backend.go .
RUN go mod init postdata_backend
RUN go mod tidy
RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o postdata_backend .

FROM alpine
WORKDIR /
COPY --from=builder /app/postdata_backend .
CMD ["./postdata_backend"]
