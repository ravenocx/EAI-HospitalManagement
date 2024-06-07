FROM golang:1.20-alpine AS builder

RUN apk update && apk add --no-cache git

# Move to working directory (/build).
WORKDIR /build

# Copy the code into the container.
COPY . .

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN go mod tidy && go mod vendor

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o apiserver .

FROM scratch

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder /build/apiserver /

# Command to run when starting the container.
ENTRYPOINT ["/apiserver"]