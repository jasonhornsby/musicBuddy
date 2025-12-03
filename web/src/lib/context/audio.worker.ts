import "$lib/context/wasm_exec.js"

let isReady = false;

async function init() {
    const go = new Go();
    const response = await fetch('/main.wasm');
    const buffer = await response.arrayBuffer();
    const { instance } = await WebAssembly.instantiate(buffer, go.importObject);
    go.run(instance);
    isReady = true;
    self.postMessage({ type: 'ready' });
}

self.onmessage = async (event) => {
    const { type, payload, id } = event.data;

    if (type === 'loadAudio') {
        let success = await loadAudio(payload);
        postMessage({
            type: 'loadAudioResult',
            payload: success,
            id
        });
    } else if (type === 'unloadAudio') {
        let success = await unloadAudio();
        postMessage({
            type: 'unloadAudioResult',
            payload: success,
            id
        });
    } else if (type === 'getAudioMetadata') {
        let success = await getAudioMetadata();
        postMessage({
            type: 'getAudioMetadataResult',
            payload: success,
            id
        });
    } else {
        console.error('Unknown message type:', type);
    }
}

init();
