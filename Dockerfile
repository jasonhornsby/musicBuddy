# Stage 1: Build Go WASM binary
FROM golang:1.25.4-alpine AS go-builder

WORKDIR /build

# Copy Go workspace and module definitions
COPY go.work go.work.sum ./
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy Go source code
COPY pkg/ ./pkg/
COPY cmd/wasm/ ./cmd/wasm/

# Build the WASM binary
# GOOS=js GOARCH=wasm tells the Go compiler to target WebAssembly
RUN GOOS=js GOARCH=wasm go build -o wasm.wasm ./cmd/wasm/main.go


# Stage 2: Build SvelteKit application
FROM node:20-alpine AS node-builder

WORKDIR /app

# Copy package files for dependency installation
COPY web/package*.json ./

# Install all dependencies (including devDependencies for the build)
RUN npm ci

# Copy the SvelteKit source code
COPY web/ .

# Copy the compiled WASM binary from the Go builder stage
# It is placed in the static folder so SvelteKit can serve it
COPY --from=go-builder /build/wasm.wasm ./static/main.wasm

# Build the SvelteKit application
RUN npm run build


# Stage 3: Production environment
FROM node:20-alpine

WORKDIR /app

# Copy the built application from the builder stage
COPY --from=node-builder /app/package*.json ./
COPY --from=node-builder /app/build ./build
COPY --from=node-builder /app/node_modules ./node_modules

# Set environment to production
ENV NODE_ENV=production
ENV PORT=3000

# Expose the port the server listens on
EXPOSE 3000

# Start the Node.js server
CMD ["node", "build"]

