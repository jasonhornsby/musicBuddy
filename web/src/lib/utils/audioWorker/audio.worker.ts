import "$lib/context/wasm_exec.js"

declare global {
    class Go {
        importObject: WebAssembly.Imports;
        run(instance: WebAssembly.Instance): Promise<void>;
        env: Record<string, string>;
    }
    interface Window {
        allocateVizBuffer: (id: string, size: number) => Uint8Array;
    }
}


console.log('[Worker] Starting wasm initialization...');

const go = new Go();
go.env = Object.assign({ GODEBUG: "gctrace=1" }, go.env);
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

self.allocateVizBuffer = (id: string, size: number) => {
    console.log('[Worker] Allocating viz buffer for id: ', id, ' size: ', size);
    const buffer = new SharedArrayBuffer(size);

    self.postMessage({
        type: 'buffer_allocated',
        id,
        buffer: buffer
    });
    return new Uint8Array(buffer)
}

self.onmessage = (event: MessageEvent) => {
    console.log('[Worker] Received message:', event.data);
    if (!goReady) {
        console.warn('[Worker] Go runtime not ready yet');
        return;
    }
}

