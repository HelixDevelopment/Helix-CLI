#!/bin/bash

# HelixCode Website Stop Script
# This script stops the website and cleans up resources

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

# Function to check if website is running
is_website_running() {
    if docker-compose ps | grep -q "Up"; then
        return 0
    else
        return 1
    fi
}

# Function to get current port
get_current_port() {
    if [ -f .website-port ]; then
        cat .website-port 2>/dev/null || echo "8000"
    else
        echo "8000"
    fi
}

# Main script
main() {
    log "Stopping HelixCode Website..."
    
    # Check if website is running
    if ! is_website_running; then
        warn "Website is not currently running"
        exit 0
    fi
    
    # Get current port for display
    local current_port=$(get_current_port)
    
    # Stop the containers
    log "Stopping Docker containers..."
    if docker-compose down; then
        log "Docker containers stopped successfully"
    else
        error "Failed to stop Docker containers"
        docker-compose logs
        exit 1
    fi
    
    # Clean up any dangling containers/networks
    log "Cleaning up resources..."
    
    # Remove dangling containers
    local dangling_containers=$(docker ps -aq -f status=exited -f name=helixcode-website)
    if [ -n "$dangling_containers" ]; then
        docker rm $dangling_containers 2>/dev/null || true
    fi
    
    # Remove dangling networks
    local dangling_networks=$(docker network ls -q -f name=helixcode-website)
    if [ -n "$dangling_networks" ]; then
        docker network rm $dangling_networks 2>/dev/null || true
    fi
    
    # Remove port file
    if [ -f .website-port ]; then
        rm .website-port
    fi
    
    # Display success message
    echo
    log "========================================"
    log "🛑 HelixCode Website Stopped Successfully!"
    log "========================================"
    info "Website was running on port: $current_port"
    info "All resources have been cleaned up"
    log "========================================"
    echo
    info "To start the website again, run: ./start-website.sh"
    echo
}

# Handle cleanup on script exit
cleanup() {
    # Remove any temporary files
    if [ -f .website-port.tmp ]; then
        rm -f .website-port.tmp
    fi
}

trap cleanup EXIT

# Run main function
main "$@"