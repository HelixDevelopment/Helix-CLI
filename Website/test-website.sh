#!/bin/bash

# HelixCode Website Test Script
# This script runs comprehensive tests on the website

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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to run HTML validation
validate_html() {
    log "Validating HTML structure..."
    
    if command_exists html5validator; then
        if html5validator --root . --ignore "node_modules/*" --ignore "*.log"; then
            log "✓ HTML validation passed"
        else
            warn "HTML validation found issues"
        fi
    else
        warn "html5validator not installed, skipping HTML validation"
        info "Install with: pip install html5validator"
    fi
}

# Function to run CSS validation
validate_css() {
    log "Validating CSS..."
    
    if command_exists css-validator; then
        for css_file in styles/*.css; do
            if [ -f "$css_file" ]; then
                if css-validator "$css_file"; then
                    log "✓ $css_file validation passed"
                else
                    warn "$css_file validation found issues"
                fi
            fi
        done
    else
        warn "css-validator not installed, skipping CSS validation"
        info "Install with: npm install -g css-validator"
    fi
}

# Function to run JavaScript syntax check
validate_javascript() {
    log "Validating JavaScript syntax..."
    
    if command_exists node; then
        for js_file in js/*.js; do
            if [ -f "$js_file" ]; then
                if node -c "$js_file"; then
                    log "✓ $js_file syntax is valid"
                else
                    error "$js_file syntax is invalid"
                    return 1
                fi
            fi
        done
    else
        warn "Node.js not installed, skipping JavaScript validation"
    fi
}

# Function to test Docker build
test_docker_build() {
    log "Testing Docker build..."
    
    if command_exists docker; then
        if docker build -t helixcode-website-test .; then
            log "✓ Docker build successful"
            
            # Test Docker run
            log "Testing Docker container..."
            if docker run -d --name helixcode-test -p 8080:80 helixcode-website-test; then
                sleep 5
                
                # Test HTTP response
                if curl -f http://localhost:8080/ > /dev/null 2>&1; then
                    log "✓ Website is accessible in Docker container"
                else
                    error "Website is not accessible in Docker container"
                    docker logs helixcode-test
                    return 1
                fi
                
                # Test health endpoint
                if curl -f http://localhost:8080/health > /dev/null 2>&1; then
                    log "✓ Health endpoint is working"
                else
                    error "Health endpoint is not working"
                    return 1
                fi
                
                # Cleanup test container
                docker stop helixcode-test >/dev/null 2>&1 || true
                docker rm helixcode-test >/dev/null 2>&1 || true
                docker rmi helixcode-website-test >/dev/null 2>&1 || true
                
            else
                error "Failed to run Docker container"
                return 1
            fi
            
        else
            error "Docker build failed"
            return 1
        fi
    else
        warn "Docker not installed, skipping Docker tests"
    fi
}

# Function to test website functionality
test_website_functionality() {
    log "Testing website functionality..."
    
    # Start website temporarily
    if ./start-website.sh > /dev/null 2>&1; then
        sleep 10
        
        local port=$(cat .website-port 2>/dev/null || echo "8000")
        
        # Test main page
        if curl -f "http://localhost:$port/" > /dev/null 2>&1; then
            log "✓ Main page loads successfully"
        else
            error "Main page failed to load"
            return 1
        fi
        
        # Test CSS loading
        if curl -f "http://localhost:$port/styles/main.css" > /dev/null 2>&1; then
            log "✓ CSS files are accessible"
        else
            error "CSS files are not accessible"
            return 1
        fi
        
        # Test JavaScript loading
        if curl -f "http://localhost:$port/js/main.js" > /dev/null 2>&1; then
            log "✓ JavaScript files are accessible"
        else
            error "JavaScript files are not accessible"
            return 1
        fi
        
        # Test fractal system
        if curl -f "http://localhost:$port/js/fractal.js" > /dev/null 2>&1; then
            log "✓ Fractal system is accessible"
        else
            error "Fractal system is not accessible"
            return 1
        fi
        
        # Stop website
        ./stop-website.sh > /dev/null 2>&1
        
    else
        error "Failed to start website for testing"
        return 1
    fi
}

# Function to test responsive design
test_responsive_design() {
    log "Testing responsive design requirements..."
    
    # Check viewport meta tag
    if grep -q '<meta name="viewport"' index.html; then
        log "✓ Viewport meta tag present"
    else
        error "Viewport meta tag missing"
        return 1
    fi
    
    # Check responsive CSS
    if grep -q '@media' styles/main.css; then
        log "✓ Responsive media queries present"
    else
        warn "No responsive media queries found"
    fi
    
    # Check mobile navigation
    if grep -q 'nav-toggle' index.html && grep -q 'nav-menu' index.html; then
        log "✓ Mobile navigation structure present"
    else
        warn "Mobile navigation structure might be incomplete"
    fi
}

# Function to test performance
test_performance() {
    log "Testing performance optimizations..."
    
    # Check for image optimization
    if [ -f "assets/logo.png" ]; then
        local logo_size=$(stat -f%z "assets/logo.png" 2>/dev/null || stat -c%s "assets/logo.png" 2>/dev/null)
        if [ "$logo_size" -lt 50000 ]; then
            log "✓ Logo image is optimized ($logo_size bytes)"
        else
            warn "Logo image might be too large ($logo_size bytes)"
        fi
    fi
    
    # Check for CSS minification
    if grep -q 'minify' package.json || [ -f "styles/main.min.css" ]; then
        log "✓ CSS minification configured"
    else
        info "Consider adding CSS minification for production"
    fi
    
    # Check for JavaScript optimization
    if grep -q 'minify' package.json || [ -f "js/main.min.js" ]; then
        log "✓ JavaScript optimization configured"
    else
        info "Consider adding JavaScript minification for production"
    fi
}

# Function to run all tests
run_all_tests() {
    local failed_tests=0
    
    log "Starting comprehensive website tests..."
    echo
    
    # Run validation tests
    validate_html || ((failed_tests++))
    echo
    
    validate_css || ((failed_tests++))
    echo
    
    validate_javascript || ((failed_tests++))
    echo
    
    # Run functionality tests
    test_docker_build || ((failed_tests++))
    echo
    
    test_website_functionality || ((failed_tests++))
    echo
    
    test_responsive_design || ((failed_tests++))
    echo
    
    test_performance || ((failed_tests++))
    echo
    
    # Summary
    if [ $failed_tests -eq 0 ]; then
        log "🎉 All tests passed successfully!"
        return 0
    else
        error "$failed_tests test(s) failed"
        return 1
    fi
}

# Main script
main() {
    case "${1:-all}" in
        "html")
            validate_html
            ;;
        "css")
            validate_css
            ;;
        "js")
            validate_javascript
            ;;
        "docker")
            test_docker_build
            ;;
        "functionality")
            test_website_functionality
            ;;
        "responsive")
            test_responsive_design
            ;;
        "performance")
            test_performance
            ;;
        "all")
            run_all_tests
            ;;
        *)
            error "Unknown test category: $1"
            echo "Available categories: html, css, js, docker, functionality, responsive, performance, all"
            exit 1
            ;;
    esac
}

# Handle signals
trap 'error "Testing interrupted"; exit 1' INT TERM

# Run main function
main "$@"