// Performance-Optimized Fractal Background System

class PerformanceFractalBackground {
    constructor() {
        this.canvas = null;
        this.ctx = null;
        this.animationId = null;
        this.isAnimating = false;
        this.currentSection = 'hero';
        this.time = 0;
        this.frameCount = 0;
        this.maxFPS = 30; // Limit FPS to reduce CPU usage
        this.lastFrameTime = 0;
        this.init();
    }

    init() {
        this.createCanvas();
        this.startAnimation();
        this.setupIntersectionObserver();
        
        // Performance monitoring
        this.setupPerformanceMonitoring();
    }

    createCanvas() {
        // Create main fractal canvas
        this.canvas = document.createElement('canvas');
        this.canvas.className = 'performance-fractal-canvas hero';
        this.canvas.style.position = 'fixed';
        this.canvas.style.top = '0';
        this.canvas.style.left = '0';
        this.canvas.style.width = '100%';
        this.canvas.style.height = '100%';
        this.canvas.style.zIndex = '-2';
        this.canvas.style.pointerEvents = 'none';
        
        // Append to body
        document.body.appendChild(this.canvas);
        
        this.ctx = this.canvas.getContext('2d');
        this.resizeCanvas();
        
        // Handle window resize
        window.addEventListener('resize', () => {
            this.resizeCanvas();
        });
    }

    resizeCanvas() {
        // Use smaller resolution for better performance
        const scale = 0.5; // 50% resolution for performance
        this.canvas.width = window.innerWidth * scale;
        this.canvas.height = window.innerHeight * scale;
        this.canvas.style.width = '100%';
        this.canvas.style.height = '100%';
    }

    draw() {
        if (!this.ctx) return;

        const currentTime = performance.now();
        const deltaTime = currentTime - this.lastFrameTime;
        
        // Limit FPS
        if (deltaTime < 1000 / this.maxFPS) {
            return;
        }
        
        this.lastFrameTime = currentTime;

        const width = this.canvas.width;
        const height = this.canvas.height;
        
        // Clear canvas with fade effect
        this.ctx.fillStyle = 'rgba(0, 0, 0, 0.1)';
        this.ctx.fillRect(0, 0, width, height);
        
        // Draw simple gradient background (much more performant)
        this.drawGradientBackground();
        
        // Draw minimal particles (limited quantity)
        this.drawParticles();
        
        this.time += deltaTime;
        this.frameCount++;
    }

    drawGradientBackground() {
        const width = this.canvas.width;
        const height = this.canvas.height;
        const time = this.time * 0.001;
        
        // Get theme colors
        const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
        
        // Create gradient
        const gradient = this.ctx.createRadialGradient(
            width * (0.3 + Math.sin(time * 0.1) * 0.1),
            height * (0.7 + Math.cos(time * 0.08) * 0.1),
            0,
            width * (0.5 + Math.sin(time * 0.05) * 0.2),
            height * (0.5 + Math.cos(time * 0.06) * 0.2),
            Math.max(width, height) * 0.8
        );
        
        if (isDark) {
            gradient.addColorStop(0, 'rgba(14, 165, 233, 0.08)');
            gradient.addColorStop(0.5, 'rgba(139, 92, 246, 0.05)');
            gradient.addColorStop(1, 'rgba(16, 185, 129, 0.03)');
        } else {
            gradient.addColorStop(0, 'rgba(14, 165, 233, 0.05)');
            gradient.addColorStop(0.5, 'rgba(139, 92, 246, 0.03)');
            gradient.addColorStop(1, 'rgba(16, 185, 129, 0.02)');
        }
        
        this.ctx.fillStyle = gradient;
        this.ctx.fillRect(0, 0, width, height);
    }

    drawParticles() {
        const width = this.canvas.width;
        const height = this.canvas.height;
        const time = this.time * 0.001;
        
        // Very limited number of particles for performance
        const particleCount = 15;
        
        // Get theme colors
        const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
        const colors = isDark ? 
            ['rgba(14, 165, 233, 0.6)', 'rgba(139, 92, 246, 0.4)', 'rgba(16, 185, 129, 0.5)'] :
            ['rgba(14, 165, 233, 0.4)', 'rgba(139, 92, 246, 0.3)', 'rgba(16, 185, 129, 0.3)'];

        for (let i = 0; i < particleCount; i++) {
            const angle = time * 0.5 + (i * Math.PI * 2) / particleCount;
            const radius = Math.sin(time * 0.2 + i) * 0.3 + 0.4;
            const x = width * (0.5 + Math.cos(angle) * radius);
            const y = height * (0.5 + Math.sin(angle) * radius);
            const size = Math.sin(time * 0.3 + i) * 2 + 3;
            const color = colors[i % colors.length];
            
            this.ctx.fillStyle = color;
            this.ctx.beginPath();
            this.ctx.arc(x, y, size, 0, Math.PI * 2);
            this.ctx.fill();
        }
    }

    startAnimation() {
        if (this.isAnimating) return;
        
        this.isAnimating = true;
        const animate = () => {
            this.draw();
            this.animationId = requestAnimationFrame(animate);
        };
        
        animate();
    }

    stopAnimation() {
        if (this.animationId) {
            cancelAnimationFrame(this.animationId);
            this.isAnimating = false;
        }
    }

    setupIntersectionObserver() {
        const sections = document.querySelectorAll('section[id]');
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    this.currentSection = entry.target.id;
                    this.updateCanvasClass();
                }
            });
        }, { threshold: 0.3 });

        sections.forEach(section => {
            observer.observe(section);
        });
    }

    updateCanvasClass() {
        if (this.canvas) {
            this.canvas.className = `performance-fractal-canvas ${this.currentSection}`;
        }
    }

    setupPerformanceMonitoring() {
        // Log performance every 10 seconds
        setInterval(() => {
            if (this.frameCount > 0) {
                const fps = Math.round((this.frameCount * 1000) / 10000);
                console.log(`🎯 Fractal Performance: ${fps} FPS`);
                this.frameCount = 0;
                
                // Auto-adjust FPS based on performance
                if (fps < 20) {
                    this.maxFPS = Math.max(15, this.maxFPS - 5);
                    console.log(`📉 Performance low, reducing FPS to ${this.maxFPS}`);
                } else if (fps > 35) {
                    this.maxFPS = Math.min(60, this.maxFPS + 5);
                    console.log(`📈 Performance good, increasing FPS to ${this.maxFPS}`);
                }
            }
        }, 10000);
    }

    destroy() {
        this.stopAnimation();
        if (this.canvas && this.canvas.parentElement) {
            this.canvas.parentElement.removeChild(this.canvas);
        }
    }
}

// Initialize performance-optimized fractal system when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    // Only initialize if device can handle it
    const isLowPerformanceDevice = 
        navigator.hardwareConcurrency < 4 || 
        !navigator.gpu ||
        window.innerWidth < 768;
    
    if (!isLowPerformanceDevice) {
        window.performanceFractalBackground = new PerformanceFractalBackground();
    } else {
        console.log('📱 Using lightweight mode for better performance');
    }
});