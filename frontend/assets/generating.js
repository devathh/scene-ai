const API_BASE = '/api/v1';
let currentToken = localStorage.getItem('token');
let ws = null;
let scenarioId = null;
let receivedSceneIds = new Set();
let isFinished = false;

const STATUS_GENERATED = 1;
const STATUS_GENERATING = 2; 
const STATUS_MODIFIED = 3;

function getScenarioIdFromUrl() {
    const params = new URLSearchParams(window.location.search);
    return params.get('id');
}

function init() {
    if (!currentToken) {
        window.location.href = '/';
        return;
    }

    scenarioId = getScenarioIdFromUrl();
    if (!scenarioId) {
        const statusText = document.getElementById('statusText');
        if (statusText) statusText.innerText = 'Error: No Scenario ID provided.';
        return;
    }

    checkScenarioStatus();
}

async function checkScenarioStatus() {
    try {
        const response = await fetch(`${API_BASE}/scenario/scenario/${scenarioId}`, {
            headers: { 'Authorization': currentToken }
        });
        
        if (!response.ok) {
            console.log('Scenario not in main DB, checking AI service...');
            checkAiStatusFallback();
            return;
        }

        const data = await response.json();
        const status = data.status;
        const scenes = data.scenes || [];

        if (scenes.length > 0) {
            renderScenes(scenes);
        }

        if (status === STATUS_GENERATED || status === STATUS_MODIFIED) {
            finishGenerationUI();
        } else {
            const statusText = document.getElementById('statusText');
            if (statusText) statusText.innerText = scenes.length > 0 ? 'Continuing generation...' : 'AI is dreaming up your scenes...';
            
            connectWebSocket();
        }

    } catch (err) {
        console.error(err);
        checkAiStatusFallback();
    }
}

async function checkAiStatusFallback() {
    try {
        const response = await fetch(`${API_BASE}/ai/scenario/${scenarioId}`, {
            headers: { 'Authorization': currentToken }
        });
        if (response.ok) {
            const data = await response.json();
            const statusText = document.getElementById('statusText');
            if (statusText) statusText.innerText = 'Connecting to AI engine...';
            connectWebSocket();
            
            const scenesResp = await fetch(`${API_BASE}/ai/scenario/${scenarioId}/scenes`, {
                headers: { 'Authorization': currentToken }
            });
            if (scenesResp.ok) {
                const scenesData = await scenesResp.json();
                if (scenesData.scenes && scenesData.scenes.length > 0) {
                    renderScenes(scenesData.scenes);
                }
            }
        } else {
            const statusText = document.getElementById('statusText');
            if (statusText) statusText.innerText = 'Scenario not found.';
        }
    } catch (e) {
        console.error('Fallback failed', e);
    }
}

function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/v1/ai/scenes/${scenarioId}`;

    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log('WebSocket connected');
    };

    ws.onmessage = (event) => {
        try {
            const scene = JSON.parse(event.data);
            addSceneToStream(scene);
        } catch (e) {
            console.error('Failed to parse WS message', e);
        }
    };

    ws.onclose = () => {
        console.log('WebSocket closed');
        finishGenerationUI();
    };

    ws.onerror = (error) => {
        console.error('WebSocket error', error);
        const statusText = document.getElementById('statusText');
        if (statusText) statusText.innerText = 'Connection closed.';
    };
}

function finishGenerationUI() {
    isFinished = true;
    if (ws) ws.close();
    
    const loader = document.getElementById('loaderSection');
    const stream = document.getElementById('streamSection');
    const statusText = document.getElementById('statusText');

    if (loader) loader.style.display = 'none';
    if (stream) stream.style.display = 'block';
    if (statusText) statusText.innerText = 'Generation Complete!';
}

function addSceneToStream(scene) {
    if (receivedSceneIds.has(scene.id)) return;
    receivedSceneIds.add(scene.id);

    const list = document.getElementById('scenesList');
    const loader = document.getElementById('loaderSection');
    const stream = document.getElementById('streamSection');

    if (loader && stream && list.children.length === 0) {
        loader.style.display = 'none';
        stream.style.display = 'block';
    }

    if (!list) return;

    const div = document.createElement('div');
    div.className = 'scene-item new';
    div.id = `scene-${scene.id}`;
    
    const order = scene.order !== undefined ? scene.order : '?';
    const title = scene.title || 'Untitled';
    const prompt = scene.video_prompt || 'Generating details...';

    div.innerHTML = `
        <div class="scene-title">Scene ${order}: ${title}</div>
        <div class="scene-desc">${prompt}</div>
    `;
    
    list.appendChild(div);
    window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });

    setTimeout(() => {
        if (div) div.classList.remove('new');
    }, 2000);
}

function renderScenes(scenes) {
    const list = document.getElementById('scenesList');
    if (!list) return;

    if (receivedSceneIds.size === 0) {
        list.innerHTML = '';
        const loader = document.getElementById('loaderSection');
        const stream = document.getElementById('streamSection');
        if (loader) loader.style.display = 'none';
        if (stream) stream.style.display = 'block';
    }

    scenes.sort((a, b) => (a.order || 0) - (b.order || 0));

    scenes.forEach(scene => {
        if (receivedSceneIds.has(scene.id)) return;
        receivedSceneIds.add(scene.id);
        
        const div = document.createElement('div');
        div.id = `scene-${scene.id}`;
        div.className = 'scene-item';
        
        const order = scene.order !== undefined ? scene.order : '?';
        const title = scene.title || 'Untitled';
        const prompt = scene.video_prompt || 'No description';

        div.innerHTML = `
            <div class="scene-title">Scene ${order}: ${title}</div>
            <div class="scene-desc">${prompt}</div>
        `;
        list.appendChild(div);
    });
}

init();