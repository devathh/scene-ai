const API_BASE = '/api/v1';
let currentToken = localStorage.getItem('token');

const STATUS_MAP = {
    1: { label: 'Generated', class: 'status-generated' },
    2: { label: 'Generating', class: 'status-generating' },
    3: { label: 'Modified', class: 'status-modified' }
};

document.addEventListener('DOMContentLoaded', () => {
    if (!currentToken) {
        window.location.href = '/';
        return;
    }
    
    updateHeader();
    loadScenarios();
});

function updateHeader() {
    const authButtons = document.getElementById('authButtons');
    const userProfile = document.getElementById('userProfile');
    const userName = document.getElementById('userName');

    if (authButtons) authButtons.classList.add('hidden');
    if (userProfile) userProfile.classList.remove('hidden');
    if (userName) userName.innerText = "User";
}

function logout() {
    localStorage.removeItem('token');
    currentToken = null;
    window.location.href = '/';
}

async function loadScenarios() {
    const container = document.getElementById('scenariosContainer');
    if (!container) return;

    try {
        const response = await fetch(`${API_BASE}/scenario/scenarios?limit=50`, {
            headers: {
                'Authorization': currentToken
            }
        });

        if (!response.ok) {
            if (response.status === 401) {
                logout();
                return;
            }
            throw new Error('Failed to load scenarios');
        }

        const data = await response.json();
        const scenarios = Array.isArray(data) ? data : (data.scenarios || []);

        renderScenarios(scenarios, container);

    } catch (err) {
        console.error(err);
        container.innerHTML = `
            <div class="empty-state">
                <p>Error loading scenarios: ${err.message}</p>
                <button onclick="location.reload()">Retry</button>
            </div>
        `;
    }
}

function renderScenarios(scenarios, container) {
    if (!scenarios || scenarios.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <h2>No scenarios yet</h2>
                <p>Create your first AI-generated story!</p>
                <button onclick="window.location.href='/'">Create New</button>
            </div>
        `;
        return;
    }

    container.innerHTML = '';
    
    scenarios.sort((a, b) => b.created_at - a.created_at);

    scenarios.forEach(scenario => {
        const statusInfo = STATUS_MAP[scenario.status] || { label: 'Unknown', class: '' };
        const date = new Date(scenario.created_at).toLocaleDateString();
        
        const card = document.createElement('div');
        card.className = 'scenario-card';
        card.onclick = () => viewScenario(scenario.id);

        card.innerHTML = `
            <div>
                <h3>${escapeHtml(scenario.title || 'Untitled Scenario')}</h3>
                <div class="prompt-preview">${escapeHtml(scenario.scenario_prompt)}</div>
            </div>
            <div class="scenario-meta">
                <span>${date}</span>
                <span class="status-badge ${statusInfo.class}">${statusInfo.label}</span>
            </div>
        `;

        container.appendChild(card);
    });
}

function viewScenario(id) {
    window.location.href = `/generating.html?id=${id}`;
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}