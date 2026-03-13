# ---------- Build Stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# เปิดใช้ Go modules ให้แน่ใจว่าใช้ go.mod ไม่ไปหาใน GOPATH/std
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server/main.go


# ---------- Run Stage ----------
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
