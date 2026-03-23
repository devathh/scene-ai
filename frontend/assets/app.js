const API_BASE = '/api/v1';
let currentToken = localStorage.getItem('token');
let currentMode = 'login';

document.addEventListener('DOMContentLoaded', () => {
    const authButtons = document.getElementById('authButtons');
    const userProfile = document.getElementById('userProfile');
    
    if (currentToken && authButtons && userProfile) {
        showLoggedInState();
    }
});

function openModal(mode) {
    currentMode = mode;
    const modal = document.getElementById('authModal');
    const title = document.getElementById('modalTitle');
    const emailInput = document.getElementById('email');
    const passwordInput = document.getElementById('password');
    const errorMsg = document.getElementById('authErrorMsg');

    if (modal) modal.style.display = 'flex';
    if (title) title.innerText = mode === 'login' ? 'Log In' : 'Sign Up';
    if (emailInput) emailInput.value = '';
    if (passwordInput) passwordInput.value = '';
    if (errorMsg) errorMsg.innerText = '';
}

function closeModal() {
    const modal = document.getElementById('authModal');
    if (modal) modal.style.display = 'none';
}

async function handleAuth(e) {
    e.preventDefault();
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const firstname = document.getElementById('firstname');
    const lastname = document.getElementById('lastname');
    
    const btn = e.target.querySelector('button');
    const originalText = btn.innerText;
    const errorMsg = document.getElementById('authErrorMsg');
    
    btn.disabled = true;
    btn.innerText = 'Processing...';
    if (errorMsg) errorMsg.innerText = '';

    try {
        let body = { email, password };
        
        if (currentMode === 'register') {
            body.firstname = firstname.value || 'User';
            body.lastname = lastname.value || 'Name';
        }

        const endpoint = currentMode === 'login' ? '/auth/login' : '/auth/register';
        const response = await fetch(`${API_BASE}${endpoint}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        });

        const data = await response.json();

        if (!response.ok) throw new Error(data.error || 'Authentication failed');

        if (currentMode === 'login') {
            currentToken = data.access;
            localStorage.setItem('token', currentToken);
            showLoggedInState();
            closeModal();
        } else {
            alert('Registration successful! Please log in.');
            openModal('login');
        }
    } catch (err) {
        if (errorMsg) errorMsg.innerText = err.message;
    } finally {
        btn.disabled = false;
        btn.innerText = originalText;
    }
}

function logout() {
    localStorage.removeItem('token');
    currentToken = null;
    location.reload();
}

function showLoggedInState() {
    const authButtons = document.getElementById('authButtons');
    const userProfile = document.getElementById('userProfile');
    const userName = document.getElementById('userName');

    if (authButtons) authButtons.classList.add('hidden');
    if (userProfile) userProfile.classList.remove('hidden');
    if (userName) userName.innerText = "User";
}

async function generateScenario() {
    if (!currentToken) {
        openModal('login');
        return;
    }

    const prompt = document.getElementById('prompt').value;
    const btn = document.getElementById('generateBtn');
    const errorDiv = document.getElementById('errorMsg');
    
    if (!prompt) {
        if (errorDiv) errorDiv.innerText = "Please enter a scenario prompt.";
        return;
    }

    if (errorDiv) errorDiv.innerText = '';
    if (btn) {
        btn.disabled = true;
        btn.innerText = 'Generating...';
    }

    try {
        const response = await fetch(`${API_BASE}/ai/scenario`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': currentToken
            },
            body: JSON.stringify({ prompt: prompt })
        });

        const data = await response.json();

        if (!response.ok) throw new Error(data.error || 'Generation failed');

        window.location.href = `/generating.html?id=${data.id}`;

    } catch (err) {
        if (errorDiv) errorDiv.innerText = err.message;
        if (btn) {
            btn.disabled = false;
            btn.innerText = 'Generate Scenario';
        }
    }
}

window.onclick = function(event) {
    const modal = document.getElementById('authModal');
    if (event.target == modal) {
        closeModal();
    }
}