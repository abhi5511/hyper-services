# Pehle yahan golang:1.21-alpine tha, isey 1.26 kar do
FROM golang:1.26-alpine AS builder 
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o hyper-id cmd/server/main.go

# Baki sab same rahega...
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/hyper-id .
COPY --from=builder /app/frontend ./frontend
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./hyper-id"]