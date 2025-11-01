#!/bin/bash

# HelixCode Website Startup Script
# This script starts the website using Docker Compose

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Function to log messages
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING:${NC} $1"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1"
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO:${NC} $1"
}

# Function to check if port is available
check_port() {
    local port=$1
    if command -v netstat >/dev/null 2>&1; then
        if netstat -tuln | grep ":$port " >/dev/null; then
            return 1
        fi
    elif command -v ss >/dev/null 2>&1; then
        if ss -tuln | grep ":$port " >/dev/null; then
            return 1
        fi
    elif command -v lsof >/dev/null 2>&1; then
        if lsof -i :$port >/dev/null 2>&1; then
            return 1
        fi
    fi
    return 0
}

# Function to find next available port
find_available_port() {
    local start_port=${1:-8000}
    local port=$start_port
    
    while ! check_port $port; do
        warn "Port $port is already in use"
        port=$((port + 1))
        if [ $port -gt 65535 ]; then
            error "No available ports found in range $start_port-65535"
            exit 1
        fi
    done
    
    echo $port
}

# Function to build latest version
build_latest_version() {
    log "Building latest version of HelixCode website..."
    
    # Check if Docker is running
    if ! docker info >/dev/null 2>&1; then
        error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    
    # Build the Docker image
    log "Building Docker image..."
    if docker-compose build --no-cache; then
        log "Docker image built successfully"
    else
        error "Failed to build Docker image"
        exit 1
    fi
}

# Function to get local IP address
get_local_ip() {
    local ip
    if command -v ip >/dev/null 2>&1; then
        ip=$(ip route get 1 | awk '{print $7; exit}')
    elif command -v ifconfig >/dev/null 2>&1; then
        ip=$(ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | head -n1)
    else
        ip="localhost"
    fi
    
    if [ -z "$ip" ] || [ "$ip" = "127.0.0.1" ]; then
        ip="localhost"
    fi
    
    echo "$ip"
}

# Main script
main() {
    log "Starting HelixCode Website..."
    
    # Check if website is already running
    if docker-compose ps | grep -q "Up"; then
        warn "Website is already running. Stopping it first..."
        ./stop-website.sh
        sleep 2
    fi
    
    # Build latest version
    build_latest_version
    
    # Find available port
    DEFAULT_PORT=8000
    WEBSITE_PORT=$(find_available_port $DEFAULT_PORT)
    
    if [ "$WEBSITE_PORT" != "$DEFAULT_PORT" ]; then
        warn "Default port $DEFAULT_PORT is occupied, using port $WEBSITE_PORT instead"
    fi
    
    export WEBSITE_PORT
    
    # Start the website
    log "Starting website on port $WEBSITE_PORT..."
    if docker-compose up -d; then
        # Wait for container to be healthy
        log "Waiting for website to become healthy..."
        local max_attempts=30
        local attempt=1
        
        while [ $attempt -le $max_attempts ]; do
            if docker-compose ps | grep -q "Up (healthy)"; then
                break
            fi
            
            if [ $attempt -eq $max_attempts ]; then
                warn "Website is taking longer than expected to start. Checking status..."
                docker-compose ps
                break
            fi
            
            sleep 2
            attempt=$((attempt + 1))
        done
        
        # Get local IP
        local ip=$(get_local_ip)
        
        # Display success message
        echo
        log "=========================================="
        log "🚀 HelixCode Website Started Successfully!"
        log "=========================================="
        info "Local URL:    http://localhost:$WEBSITE_PORT"
        info "Network URL:  http://$ip:$WEBSITE_PORT"
        info "Health Check: http://localhost:$WEBSITE_PORT/health"
        log "=========================================="
        echo
        info "To stop the website, run: ./stop-website.sh"
        info "To view logs, run: docker-compose logs -f"
        echo
        
        # Open browser if requested
        if [ "$1" = "--open" ] || [ "$2" = "--open" ]; then
            if command -v xdg-open >/dev/null 2>&1; then
                xdg-open "http://localhost:$WEBSITE_PORT"
            elif command -v open >/dev/null 2>&1; then
                open "http://localhost:$WEBSITE_PORT"
            else
                warn "Could not automatically open browser. Please visit the URL above."
            fi
        fi
        
    else
        error "Failed to start website"
        docker-compose logs
        exit 1
    fi
}

# Handle signals
trap 'error "Script interrupted"; exit 1' INT TERM

# Run main function
main "$@"