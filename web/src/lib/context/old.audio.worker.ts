import "$lib/context/wasm_exec.js"

import { audioActions } from "./audio.actions.ts";

let isReady = false;

async function init() {
    console.log("Initializing audio worker");
    const go = new Go();
    const response = await fetch('/main.wasm');
    const buffer = await response.arrayBuffer();
    const { instance } = await WebAssembly.instantiate(buffer, go.importObject);
    go.run(instance);
    isReady = true;
    self.postMessage({ type: 'ready' });
    console.log("Audio worker initialized");
}

self.onmessage = async (event) => {
    const { type, payload, id } = event.data;

    console.log("Received message: ", type);

    for (const action of audioActions) {
        if (type === action.requestKey) {
            const handler = self[action.actionKey as keyof Window];
            const result = await handler(payload);
            postMessage({
                type: action.responseKey,
                payload: result,
                id
            });
            return;
        }
    }
    console.error('Unknown message type:', type);
}

init();
