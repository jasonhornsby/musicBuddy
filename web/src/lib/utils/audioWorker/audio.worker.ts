import "$lib/context/wasm_exec.js"

declare global {
    class Go {
        importObject: WebAssembly.Imports;
        run(instance: WebAssembly.Instance): Promise<void>;
    }
}


console.log('[Worker] Starting wasm initialization...');

const go = new Go();
let goReady = false;

WebAssembly.instantiateStreaming(fetch('/main.wasm'), go.importObject)
    .then((result) => {
        console.log('[Worker] Wasm module loaded');
        go.run(result.instance);
        console.log('[Worker] Go runtime initialized');

        postMessage({ type: 'worker_ready', timeStamp: Date.now() });
        goReady = true;
    }).catch((error) => {
        console.error('[Worker] Error loading wasm:', error);
        postMessage({ type: 'worker_error', error: error.message, timeStamp: Date.now() });
    });

self.onmessage = (event: MessageEvent) => {
    console.log('[Worker] Received message:', event.data);
    if (!goReady) {
        console.warn('[Worker] Go runtime not ready yet');
        return;
    }
}

