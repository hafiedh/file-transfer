FROM golang:latest as builder
LABEL MAINTAINER="Muhammad Hafiedh"

WORKDIR /go/src/file-service
COPY . .

RUN go mod download && \
   go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/file-service .

# Final stage using `scratch`
FROM scratch

# Copy the built binary and environment file
COPY --from=builder /go/bin/file-service /app/file-service
COPY --from=builder /go/src/file-service/.env /app/.env

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER 1000
WORKDIR /app

EXPOSE 8090
EXPOSE 50052

ENTRYPOINT ["/app/file-service", "-env", "/app/.env"]