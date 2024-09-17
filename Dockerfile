# Builder container
FROM golang:1.22 AS builder

# Set working directory
WORKDIR /src

# Update and install apt dependencies
RUN apt-get update && apt-get install -q -y stockfish

# Copy go module files
COPY go.mod go.sum ./

# Download go module dependencies
RUN go mod download all

# Copy source files and build executable
COPY . ./
RUN CGO_ENABLED=0 go build -v -o /bin/app ./cmd/discordbot/main.go

# Production container
FROM debian:bullseye-slim AS final

# Set working directory
WORKDIR /app

# Update and install apt dependencies
RUN apt-get update && \
    apt-get install -q -y \
    ffmpeg \
    inkscape

# Copy stockfish from builder and add to path
COPY --from=builder /usr/games/stockfish /usr/games
ENV PATH /usr/games:$PATH

COPY .env ./

# Copy executable from builder
COPY --from=builder /bin/app ./

# Start app
CMD ["./app"]
