const CustomCursor = {
    cursor: null,
    follower: null,
    dot: null,
    mouseX: 0,
    mouseY: 0,
    followerX: 0,
    followerY: 0,
    isHovering: false,

    init() {
        this.cursor = document.createElement('div');
        this.follower = document.createElement('div');
        this.dot = document.createElement('div');

        this.cursor.className = 'custom-cursor';
        this.follower.className = 'custom-cursor-follower';
        this.dot.className = 'custom-cursor-dot';

        document.body.appendChild(this.cursor);
        document.body.appendChild(this.follower);
        document.body.appendChild(this.dot);

        const interactiveElements = document.querySelectorAll('button, a, .scenario-card, input, textarea, [onclick]');
        
        interactiveElements.forEach(el => {
            el.addEventListener('mouseenter', () => this.enterInteractive());
            el.addEventListener('mouseleave', () => this.leaveInteractive());
        });

        document.addEventListener('mousemove', (e) => {
            this.mouseX = e.clientX;
            this.mouseY = e.clientY;
            
            this.dot.style.left = this.mouseX + 'px';
            this.dot.style.top = this.mouseY + 'px';
        });

        this.animate();
    },

    enterInteractive() {
        this.isHovering = true;
        if (this.follower) {
            this.follower.style.width = '60px';
            this.follower.style.height = '60px';
            this.follower.style.backgroundColor = 'rgba(0, 0, 0, 0.05)';
            this.follower.style.borderColor = 'rgba(0, 0, 0, 0.1)';
        }
        if (this.cursor) {
            this.cursor.style.transform = 'translate(-50%, -50%) scale(0.5)';
        }
        if (this.dot) {
            this.dot.style.transform = 'translate(-50%, -50%) scale(1.5)';
        }
    },

    leaveInteractive() {
        this.isHovering = false;
        if (this.follower) {
            this.follower.style.width = '40px';
            this.follower.style.height = '40px';
            this.follower.style.backgroundColor = 'rgba(255, 255, 255, 0.3)';
            this.follower.style.borderColor = 'rgba(0, 0, 0, 0.2)';
        }
        if (this.cursor) {
            this.cursor.style.transform = 'translate(-50%, -50%) scale(1)';
        }
        if (this.dot) {
            this.dot.style.transform = 'translate(-50%, -50%) scale(1)';
        }
    },

    animate() {
        const speed = 0.15;

        this.followerX += (this.mouseX - this.followerX) * speed;
        this.followerY += (this.mouseY - this.followerY) * speed;

        this.cursor.style.left = this.mouseX + 'px';
        this.cursor.style.top = this.mouseY + 'px';
        
        this.follower.style.left = this.followerX + 'px';
        this.follower.style.top = this.followerY + 'px';

        requestAnimationFrame(() => this.animate());
    }
};

document.addEventListener('DOMContentLoaded', () => {
    CustomCursor.init();
});