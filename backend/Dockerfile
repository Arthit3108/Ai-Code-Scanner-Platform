FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api/server.go

FROM alpine:latest

# Install Git and scanner tools (Trivy, Gitleaks)
RUN apk add --no-cache git curl bash && \
    curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin && \
    curl -sfL https://raw.githubusercontent.com/zricethezav/gitleaks/master/install.sh | sh -s -- -b /usr/local/bin

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 3000

CMD ["./app"]