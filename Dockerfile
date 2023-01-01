FROM golang:1.19-alpine AS gobuilder

LABEL maintainer="Julian Beck <mail@julianbeck.com (https://juli.sh/)"

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY internal/ ./internal
COPY main.go ./
# Set necessary environment variables needed for our image and build the API server.
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
	go build -ldflags="-s -w" -o iplocation .

FROM alpine:latest

# Move Files from build steps into the container
COPY --from=gobuilder ["/build/iplocation", "/"]
ENTRYPOINT ["/iplocation"]