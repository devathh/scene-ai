const MagicInteraction = {
    init() {
        this.createSparklesOnHover();
    },

    createSparklesOnHover() {
        const interactiveElements = document.querySelectorAll('button, a, .scenario-card, [onclick]');
        
        interactiveElements.forEach(el => {
            el.addEventListener('mousemove', (e) => {
                this.createSparkle(e.clientX, e.clientY);
            });
        });
    },

    createSparkle(x, y) {
        const sparkle = document.createElement('div');
        sparkle.classList.add('magic-sparkle');
        sparkle.style.left = `${x}px`;
        sparkle.style.top = `${y}px`;
        
        const size = Math.random() * 6 + 4;
        sparkle.style.width = `${size}px`;
        sparkle.style.height = `${size}px`;
        
        const colors = ['#000', '#444', '#888', '#ccc'];
        sparkle.style.background = colors[Math.floor(Math.random() * colors.length)];
        
        document.body.appendChild(sparkle);

        if (typeof anime !== 'undefined') {
            const angle = Math.random() * Math.PI * 2;
            const velocity = Math.random() * 100 + 50;
            const tx = Math.cos(angle) * velocity;
            const ty = Math.sin(angle) * velocity;

            anime({
                targets: sparkle,
                translateX: tx,
                translateY: ty,
                scale: [1, 0],
                opacity: [1, 0],
                duration: 800,
                easing: 'easeOutExpo',
                complete: () => sparkle.remove()
            });
        } else {
            sparkle.style.transition = 'all 0.8s ease-out';
            sparkle.style.transform = `translate(${(Math.random()*200-100)}px, ${(Math.random()*200-100)}px) scale(0)`;
            sparkle.style.opacity = '0';
            setTimeout(() => sparkle.remove(), 800);
        }
    }
};

document.addEventListener('DOMContentLoaded', () => {
    MagicInteraction.init();
});