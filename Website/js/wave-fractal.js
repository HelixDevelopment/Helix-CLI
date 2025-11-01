// Enhanced Wave Fractal Background System

class WaveFractalBackground {
    constructor() {
        this.canvas = null;
        this.ctx = null;
        this.particles = [];
        this.waves = [];
        this.animationId = null;
        this.isAnimating = false;
        this.currentSection = 'hero';
        this.time = 0;
        this.init();
    }

    init() {
        this.createCanvas();
        this.createWaves();
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
        this.canvas.className = 'wave-fractal-canvas hero';
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
            this.draw();
        });
    }

    resizeCanvas() {
        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;
    }

    createWaves() {
        // Create multiple wave layers
        this.waves = [
            {
                amplitude: 0.8,
                frequency: 0.002,
                speed: 0.0003,
                color: 'rgba(14, 165, 233, 0.15)', // Primary blue
                blur: 15
            },
            {
                amplitude: 0.6,
                frequency: 0.003,
                speed: 0.0004,
                color: 'rgba(139, 92, 246, 0.12)', // Accent purple
                blur: 12
            },
            {
                amplitude: 0.4,
                frequency: 0.004,
                speed: 0.0005,
                color: 'rgba(16, 185, 129, 0.1)', // Accent green
                blur: 10
            },
            {
                amplitude: 0.3,
                frequency: 0.005,
                speed: 0.0006,
                color: 'rgba(245, 158, 11, 0.08)', // Accent orange
                blur: 8
            }
        ];
    }

    drawWave(wave, time) {
        const width = this.canvas.width;
        const height = this.canvas.height;
        const ctx = this.ctx;

        ctx.save();
        
        // Apply blur effect
        ctx.filter = `blur(${wave.blur}px)`;
        
        ctx.beginPath();
        
        // Start from left
        ctx.moveTo(0, height / 2);
        
        // Draw wave
        for (let x = 0; x <= width; x += 2) {
            const y = height / 2 + 
                Math.sin(x * wave.frequency + time * wave.speed) * wave.amplitude * 100 +
                Math.cos(x * wave.frequency * 0.7 + time * wave.speed * 1.3) * wave.amplitude * 50;
            
            ctx.lineTo(x, y);
        }
        
        // Complete the shape
        ctx.lineTo(width, height);
        ctx.lineTo(0, height);
        ctx.closePath();
        
        // Create gradient
        const gradient = ctx.createLinearGradient(0, height / 2, 0, height);
        gradient.addColorStop(0, wave.color);
        gradient.addColorStop(1, 'transparent');
        
        ctx.fillStyle = gradient;
        ctx.fill();
        
        ctx.restore();
    }

    drawFractalPatterns(time) {
        const width = this.canvas.width;
        const height = this.canvas.height;
        const ctx = this.ctx;

        // Get theme colors
        const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
        const colors = isDark ? 
            ['rgba(14, 165, 233, 0.08)', 'rgba(139, 92, 246, 0.06)', 'rgba(16, 185, 129, 0.05)'] :
            ['rgba(14, 165, 233, 0.06)', 'rgba(139, 92, 246, 0.04)', 'rgba(16, 185, 129, 0.03)'];

        // Draw floating fractal particles
        for (let i = 0; i < 50; i++) {
            const x = (Math.sin(time * 0.001 + i) * 0.5 + 0.5) * width;
            const y = (Math.cos(time * 0.001 + i * 0.7) * 0.5 + 0.5) * height;
            const size = Math.sin(time * 0.002 + i) * 2 + 3;
            const color = colors[i % colors.length];
            
            ctx.save();
            ctx.filter = `blur(${Math.sin(time * 0.001 + i) * 2 + 3}px)`;
            ctx.fillStyle = color;
            ctx.beginPath();
            ctx.arc(x, y, size, 0, Math.PI * 2);
            ctx.fill();
            ctx.restore();
        }

        // Draw Julia set patterns
        this.drawJuliaSet(time, colors);
    }

    drawJuliaSet(time, colors) {
        const width = this.canvas.width;
        const height = this.canvas.height;
        const ctx = this.ctx;

        const scale = 2.5 + Math.sin(time * 0.0005) * 0.5;
        const offsetX = Math.sin(time * 0.0003) * 0.3;
        const offsetY = Math.cos(time * 0.0004) * 0.3;

        for (let x = 0; x < width; x += 6) {
            for (let y = 0; y < height; y += 6) {
                const zx = (x - width / 2) / (width * 0.4) * scale + offsetX;
                const zy = (y - height / 2) / (height * 0.4) * scale + offsetY;
                
                let zx0 = zx;
                let zy0 = zy;
                let iteration = 0;
                const maxIterations = 15;

                while (zx0 * zx0 + zy0 * zy0 < 4 && iteration < maxIterations) {
                    const tmp = zx0 * zx0 - zy0 * zy0 + Math.sin(time * 0.0008) * 0.7;
                    zy0 = 2 * zx0 * zy0 + Math.cos(time * 0.0006) * 0.7;
                    zx0 = tmp;
                    iteration++;
                }

                if (iteration < maxIterations) {
                    const colorIndex = iteration % colors.length;
                    ctx.save();
                    ctx.filter = `blur(${Math.random() * 2 + 1}px)`;
                    ctx.fillStyle = colors[colorIndex];
                    ctx.beginPath();
                    ctx.arc(x, y, 1, 0, Math.PI * 2);
                    ctx.fill();
                    ctx.restore();
                }
            }
        }
    }

    draw() {
        if (!this.ctx) return;

        const width = this.canvas.width;
        const height = this.canvas.height;
        
        // Clear canvas with fade effect
        this.ctx.fillStyle = 'rgba(0, 0, 0, 0.05)';
        this.ctx.fillRect(0, 0, width, height);
        
        // Draw waves
        this.waves.forEach(wave => {
            this.drawWave(wave, this.time);
        });
        
        // Draw fractal patterns
        this.drawFractalPatterns(this.time);
        
        this.time += 16; // ~60fps
    }

    createParticles() {
        const container = document.createElement('div');
        container.className = 'wave-particles-container';
        document.body.appendChild(container);

        // Create particles
        for (let i = 0; i < 30; i++) {
            this.createParticle(container);
        }
    }

    createParticle(container) {
        const particle = document.createElement('div');
        particle.className = 'wave-particle';
        
        // Random properties
        const size = Math.random() * 4 + 1;
        const left = Math.random() * 100;
        const opacity = Math.random() * 0.3 + 0.1;
        const duration = Math.random() * 25 + 15;
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
        particle.style.filter = `blur(${Math.random() * 2 + 1}px)`;
        
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
            this.canvas.className = `wave-fractal-canvas ${this.currentSection}`;
        }
    }

    handleThemeChange() {
        // Redraw when theme changes
        const observer = new MutationObserver((mutations) => {
            mutations.forEach((mutation) => {
                if (mutation.attributeName === 'data-theme') {
                    this.draw();
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
        const particlesContainer = document.querySelector('.wave-particles-container');
        if (particlesContainer) {
            particlesContainer.remove();
        }
    }
}

// Initialize wave fractal system when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.waveFractalBackground = new WaveFractalBackground();
});