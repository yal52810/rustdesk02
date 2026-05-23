# Multi-stage build for RustDesk Web Frontend
FROM node:18-alpine AS builder

WORKDIR /build

# Copy package files
COPY package*.json ./

# Install all dependencies (including devDependencies for build)
RUN npm install

# Copy source code
COPY . .

# Build production bundle
RUN npm run build

# Runtime stage - Nginx to serve static files
FROM nginx:alpine

# Copy built files
COPY --from=builder /build/dist /usr/share/nginx/html

# Copy nginx configuration (optional, can be overridden by docker-compose)
# COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost || exit 1

# Run nginx
CMD ["nginx", "-g", "daemon off;"]
