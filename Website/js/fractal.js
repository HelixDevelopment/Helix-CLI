// Dynamic Fractal Background System

class FractalBackground {
    constructor() {
        this.canvas = null;
        this.ctx = null;
        this.particles = [];
        this.animationId = null;
        this.isAnimating = false;
        this.currentSection = 'hero';
        this.init();
    }

    init() {
        this.createCanvas();
        this.createParticles();
        this.startAnimation();
        this.setupIntersectionObserver();
        
        // Listen for theme changes
        document.addEventListener('DOMContentLoaded', () => {
            this.handleThemeChange();
        });
    }

    createCanvas() {
        // Create main fractal canvas
        this.canvas = document.createElement('canvas');
        this.canvas.className = 'fractal-canvas hero';
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
            this.drawFractal();
        });
    }

    resizeCanvas() {
        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;
    }

    drawFractal() {
        if (!this.ctx) return;

        const width = this.canvas.width;
        const height = this.canvas.height;
        
        // Clear canvas
        this.ctx.clearRect(0, 0, width, height);
        
        // Get theme colors
        const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
        const colors = isDark ? 
            ['rgba(14, 165, 233, 0.1)', 'rgba(139, 92, 246, 0.08)', 'rgba(16, 185, 129, 0.06)'] :
            ['rgba(14, 165, 233, 0.08)', 'rgba(139, 92, 246, 0.06)', 'rgba(16, 185, 129, 0.04)'];

        // Draw fractal patterns
        this.drawJuliaFractal(width, height, colors);
        this.drawMandelbrotFractal(width, height, colors);
    }

    drawJuliaFractal(width, height, colors) {
        const time = Date.now() * 0.001;
        const scale = 2.5;
        const offsetX = Math.sin(time * 0.3) * 0.2;
        const offsetY = Math.cos(time * 0.2) * 0.2;

        for (let x = 0; x < width; x += 4) {
            for (let y = 0; y < height; y += 4) {
                const zx = (x - width / 2) / (width * 0.4) * scale + offsetX;
                const zy = (y - height / 2) / (height * 0.4) * scale + offsetY;
                
                let zx0 = zx;
                let zy0 = zy;
                let iteration = 0;
                const maxIterations = 20;

                while (zx0 * zx0 + zy0 * zy0 < 4 && iteration < maxIterations) {
                    const tmp = zx0 * zx0 - zy0 * zy0 + Math.sin(time * 0.5) * 0.7;
                    zy0 = 2 * zx0 * zy0 + Math.cos(time * 0.3) * 0.7;
                    zx0 = tmp;
                    iteration++;
                }

                if (iteration < maxIterations) {
                    const colorIndex = iteration % colors.length;
                    this.ctx.fillStyle = colors[colorIndex];
                    this.ctx.fillRect(x, y, 2, 2);
                }
            }
        }
    }

    drawMandelbrotFractal(width, height, colors) {
        const time = Date.now() * 0.001;
        const scale = 3;
        const offsetX = Math.cos(time * 0.4) * 0.3;
        const offsetY = Math.sin(time * 0.35) * 0.3;

        for (let x = 0; x < width; x += 6) {
            for (let y = 0; y < height; y += 6) {
                const cx = (x - width / 2) / (width * 0.3) * scale + offsetX;
                const cy = (y - height / 2) / (height * 0.3) * scale + offsetY;
                
                let zx = 0;
                let zy = 0;
                let iteration = 0;
                const maxIterations = 15;

                while (zx * zx + zy * zy < 4 && iteration < maxIterations) {
                    const tmp = zx * zx - zy * zy + cx;
                    zy = 2 * zx * zy + cy;
                    zx = tmp;
                    iteration++;
                }

                if (iteration < maxIterations) {
                    const colorIndex = (iteration + 1) % colors.length;
                    this.ctx.fillStyle = colors[colorIndex];
                    this.ctx.beginPath();
                    this.ctx.arc(x, y, 1, 0, Math.PI * 2);
                    this.ctx.fill();
                }
            }
        }
    }

    createParticles() {
        const container = document.createElement('div');
        container.className = 'particles-container';
        document.body.appendChild(container);

        // Create particles
        for (let i = 0; i < 50; i++) {
            this.createParticle(container);
        }
    }

    createParticle(container) {
        const particle = document.createElement('div');
        particle.className = 'particle';
        
        // Random properties
        const size = Math.random() * 3 + 1;
        const left = Math.random() * 100;
        const opacity = Math.random() * 0.4 + 0.1;
        const duration = Math.random() * 20 + 10;
        const delay = Math.random() * 20;
        
        // Get theme-appropriate colors
        const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
        const colors = isDark ? 
            ['rgba(14, 165, 233, 0.8)', 'rgba(139, 92, 246, 0.6)', 'rgba(16, 185, 129, 0.7)'] :
            ['rgba(14, 165, 233, 0.6)', 'rgba(139, 92, 246, 0.4)', 'rgba(16, 185, 129, 0.5)'];
        
        const color = colors[Math.floor(Math.random() * colors.length)];
        
        particle.style.width = `${size}px`;
        particle.style.height = `${size}px`;
        particle.style.left = `${left}%`;
        particle.style.background = color;
        particle.style.opacity = opacity;
        particle.style.animationDuration = `${duration}s`;
        particle.style.animationDelay = `${delay}s`;
        
        container.appendChild(particle);
        this.particles.push(particle);

        // Remove and recreate particle after animation
        setTimeout(() => {
            if (particle.parentElement) {
                particle.remove();
                this.particles = this.particles.filter(p => p !== particle);
                this.createParticle(container);
            }
        }, (duration + delay) * 1000);
    }

    startAnimation() {
        if (this.isAnimating) return;
        
        this.isAnimating = true;
        const animate = () => {
            this.drawFractal();
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
            this.canvas.className = `fractal-canvas ${this.currentSection}`;
        }
    }

    handleThemeChange() {
        // Redraw fractal when theme changes
        const observer = new MutationObserver((mutations) => {
            mutations.forEach((mutation) => {
                if (mutation.attributeName === 'data-theme') {
                    this.drawFractal();
                    // Update particles
                    this.particles.forEach(particle => {
                        const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
                        const colors = isDark ? 
                            ['rgba(14, 165, 233, 0.8)', 'rgba(139, 92, 246, 0.6)', 'rgba(16, 185, 129, 0.7)'] :
                            ['rgba(14, 165, 233, 0.6)', 'rgba(139, 92, 246, 0.4)', 'rgba(16, 185, 129, 0.5)'];
                        const color = colors[Math.floor(Math.random() * colors.length)];
                        particle.style.background = color;
                    });
                }
            });
        });

        observer.observe(document.documentElement, {
            attributes: true,
            attributeFilter: ['data-theme']
        });
    }

    destroy() {
        this.stopAnimation();
        if (this.canvas && this.canvas.parentElement) {
            this.canvas.parentElement.removeChild(this.canvas);
        }
        const particlesContainer = document.querySelector('.particles-container');
        if (particlesContainer) {
            particlesContainer.remove();
        }
    }
}

// Initialize fractal system when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.fractalBackground = new FractalBackground();
});