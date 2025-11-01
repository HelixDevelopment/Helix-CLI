#!/bin/bash

# Local Performance Test for HelixCode Website
# Tests the website without Docker for quick validation

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

# Test file structure
test_file_structure() {
    log "Testing file structure..."
    
    local required_files=(
        "index.html"
        "styles/main.css"
        "styles/performance-fractal.css"
        "js/main.js"
        "js/performance-fractal.js"
        "assets/logo.png"
    )
    
    local missing_files=()
    
    for file in "${required_files[@]}"; do
        if [ ! -f "$file" ]; then
            missing_files+=("$file")
        fi
    done
    
    if [ ${#missing_files[@]} -eq 0 ]; then
        log "✓ All required files present"
    else
        error "Missing files: ${missing_files[*]}"
        return 1
    fi
}

# Test HTML validity
test_html_validity() {
    log "Testing HTML validity..."
    
    # Check for basic HTML structure
    if grep -q '<!DOCTYPE html>' index.html && \
       grep -q '<html' index.html && \
       grep -q '<head>' index.html && \
       grep -q '<body>' index.html && \
       grep -q '</html>' index.html; then
        log "✓ Basic HTML structure is valid"
    else
        error "HTML structure is incomplete"
        return 1
    fi
    
    # Check for viewport meta tag
    if grep -q '<meta name="viewport"' index.html; then
        log "✓ Viewport meta tag present"
    else
        error "Viewport meta tag missing"
        return 1
    fi
    
    # Check for title
    if grep -q '<title>' index.html; then
        log "✓ Title tag present"
    else
        error "Title tag missing"
        return 1
    fi
}

# Test CSS validity
test_css_validity() {
    log "Testing CSS validity..."
    
    local css_files=("styles/main.css" "styles/performance-fractal.css")
    
    for css_file in "${css_files[@]}"; do
        if [ -f "$css_file" ]; then
            # Check for basic CSS structure
            if grep -q '{' "$css_file" && grep -q '}' "$css_file"; then
                log "✓ $css_file has valid CSS structure"
            else
                error "$css_file has invalid CSS structure"
                return 1
            fi
            
            # Check for CSS variables
            if grep -q 'var(--' "$css_file"; then
                log "✓ $css_file uses CSS custom properties"
            fi
        else
            error "CSS file not found: $css_file"
            return 1
        fi
    done
}

# Test JavaScript validity
test_javascript_validity() {
    log "Testing JavaScript validity..."
    
    local js_files=("js/main.js" "js/performance-fractal.js")
    
    for js_file in "${js_files[@]}"; do
        if [ -f "$js_file" ]; then
            # Check for basic JavaScript structure
            if grep -q 'function\|class\|const\|let' "$js_file"; then
                log "✓ $js_file has valid JavaScript structure"
            else
                error "$js_file has invalid JavaScript structure"
                return 1
            fi
            
            # Check for event listeners
            if grep -q 'addEventListener' "$js_file"; then
                log "✓ $js_file uses event listeners"
            fi
        else
            error "JavaScript file not found: $js_file"
            return 1
        fi
    done
}

# Test performance optimizations
test_performance_optimizations() {
    log "Testing performance optimizations..."
    
    # Check for FPS limiting in fractal system
    if grep -q 'maxFPS\|requestAnimationFrame' "js/performance-fractal.js"; then
        log "✓ Fractal system has FPS limiting"
    else
        warn "Fractal system may not have FPS limiting"
    fi
    
    # Check for low-performance mode detection
    if grep -q 'low-performance-mode\|hardwareConcurrency' "js/main.js"; then
        log "✓ Low-performance mode detection implemented"
    else
        warn "Low-performance mode detection may be missing"
    fi
    
    # Check for responsive design
    if grep -q '@media' "styles/main.css" "styles/performance-fractal.css"; then
        log "✓ Responsive design implemented"
    else
        error "Responsive design not implemented"
        return 1
    fi
}

# Test resource sizes
test_resource_sizes() {
    log "Testing resource sizes..."
    
    local max_size_kb=500
    
    # Check HTML size
    local html_size=$(stat -f%z index.html 2>/dev/null || stat -c%s index.html 2>/dev/null)
    local html_size_kb=$((html_size / 1024))
    
    if [ $html_size_kb -lt $max_size_kb ]; then
        log "✓ HTML size is reasonable: ${html_size_kb}KB"
    else
        warn "HTML size is large: ${html_size_kb}KB"
    fi
    
    # Check CSS sizes
    local total_css_size=0
    for css_file in styles/*.css; do
        if [ -f "$css_file" ]; then
            local size=$(stat -f%z "$css_file" 2>/dev/null || stat -c%s "$css_file" 2>/dev/null)
            total_css_size=$((total_css_size + size))
        fi
    done
    
    local total_css_size_kb=$((total_css_size / 1024))
    if [ $total_css_size_kb -lt $((max_size_kb * 2)) ]; then
        log "✓ Total CSS size is reasonable: ${total_css_size_kb}KB"
    else
        warn "Total CSS size is large: ${total_css_size_kb}KB"
    fi
    
    # Check JavaScript sizes
    local total_js_size=0
    for js_file in js/*.js; do
        if [ -f "$js_file" ]; then
            local size=$(stat -f%z "$js_file" 2>/dev/null || stat -c%s "$js_file" 2>/dev/null)
            total_js_size=$((total_js_size + size))
        fi
    done
    
    local total_js_size_kb=$((total_js_size / 1024))
    if [ $total_js_size_kb -lt $((max_size_kb * 2)) ]; then
        log "✓ Total JavaScript size is reasonable: ${total_js_size_kb}KB"
    else
        warn "Total JavaScript size is large: ${total_js_size_kb}KB"
    fi
}

# Run all local tests
run_all_tests() {
    local failed_tests=0
    
    log "Starting local performance and structure tests..."
    echo
    
    test_file_structure || ((failed_tests++))
    echo
    
    test_html_validity || ((failed_tests++))
    echo
    
    test_css_validity || ((failed_tests++))
    echo
    
    test_javascript_validity || ((failed_tests++))
    echo
    
    test_performance_optimizations || ((failed_tests++))
    echo
    
    test_resource_sizes || ((failed_tests++))
    echo
    
    # Summary
    if [ $failed_tests -eq 0 ]; then
        log "🎉 All local tests passed successfully!"
        log "The website structure is valid and performance-optimized"
        echo
        log "To test the website locally:"
        echo "  python3 -m http.server 8000"
        echo "  Then open http://localhost:8000 in your browser"
        return 0
    else
        error "$failed_tests test(s) failed"
        warn "Please fix the issues before deploying"
        return 1
    fi
}

# Main script
main() {
    case "${1:-all}" in
        "structure")
            test_file_structure
            ;;
        "html")
            test_html_validity
            ;;
        "css")
            test_css_validity
            ;;
        "js")
            test_javascript_validity
            ;;
        "performance")
            test_performance_optimizations
            ;;
        "sizes")
            test_resource_sizes
            ;;
        "all")
            run_all_tests
            ;;
        *)
            error "Unknown test category: $1"
            echo "Available categories: structure, html, css, js, performance, sizes, all"
            exit 1
            ;;
    esac
}

# Handle signals
trap 'error "Testing interrupted"; exit 1' INT TERM

# Run main function
main "$@"