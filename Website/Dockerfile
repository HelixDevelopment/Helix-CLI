# Multi-stage build for HelixCode Website
FROM node:18-alpine AS builder

# Install dependencies for building
RUN apk add --no-cache \
    python3 \
    make \
    g++ \
    git

# Create app directory
WORKDIR /app

# Copy package files
COPY package.json ./

# Install dependencies
RUN npm install --omit=dev

# Copy source code
COPY . .

# Build the application (if needed)
# RUN npm run build

# Production stage
FROM nginx:alpine

# Install curl for health checks
RUN apk add --no-cache curl

# Remove default nginx website
RUN rm -rf /usr/share/nginx/html/*

# Copy built application from builder stage
COPY --from=builder /app /usr/share/nginx/html

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Create health check endpoint
RUN echo '{"status": "healthy", "service": "helixcode-website"}' > /usr/share/nginx/html/health.json

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/health.json || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]