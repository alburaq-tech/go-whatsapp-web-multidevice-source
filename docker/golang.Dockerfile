############################
# STEP 1 build executable binary
############################
FROM golang:1.26.10-alpine3.23 AS builder
RUN apk update && apk add --no-cache gcc musl-dev gcompat
WORKDIR /whatsapp
COPY ./src .

# Fetch dependencies.
RUN go mod download
# Build the binary with optimizations
RUN go build -a -ldflags="-w -s" -o /app/whatsapp

#############################
## STEP 2 build a smaller image
#############################
FROM alpine:3.23
# apk upgrade: take patched Alpine packages (brotli-libs >= 1.2.1-r0) even when older layers are cached
RUN apk upgrade --no-cache && apk add --no-cache ffmpeg libwebp-tools tzdata
ENV TZ=UTC
WORKDIR /app
# Copy compiled from builder.
COPY --from=builder /app/whatsapp /app/whatsapp
# Run the binary.
ENTRYPOINT ["/app/whatsapp"]

CMD [ "rest" ]