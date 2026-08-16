# WebServe

## Description

WebServe is a simple Webserver for website. if you have a website to host WebServe is probabbly capable of running the
website efficiently and packaged into one Binary.

## Usage

we imagine we have an astro web project inside the "astro-web" directory

```go
//go:embed astro-web/dist/*
var embeddedContent embed.FS

func main() {
webFS, err := fs.Sub(embeddedContent, "astro-web/dist")
if err != nil {
fmt.Printf("Error Creating Sub FS: %s", err.Error())
return
}

r := server.New(80, webFS)

err = r.Run()
if err != nil {
fmt.Printf("error: %s", err.Error())
return
}
}
```

## Example Docker File

```yaml
# stage 1: build frontend to single index.html
FROM oven/bun AS frontend

WORKDIR /app

  # Copy package.json and bun.lock
COPY astro-web/package.json astro-web/bun.lock ./

  # Install dependencies#
RUN bun install

  # Copy the frontend
COPY astro-web /app/

  # Build frontend
RUN bun run build

  # stage 2: Build GO binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./

RUN go mod download

  # Copy frontend build files
COPY --from=frontend /app/dist /app/astro-web/dist

  # Copy go source files
COPY *.go .
COPY internal ./internal

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main .

  # stage 3: build minimal run image
FROM scratch

WORKDIR /app

  # Copy CA certificates from builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

  # Copy GO Binary
COPY --from=builder /app/main .

  # Expose Port
EXPOSE 80

  # Command to run Application
ENTRYPOINT ["./main"]

```