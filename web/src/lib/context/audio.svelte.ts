import { getContext, setContext } from "svelte";
import { toast } from "svelte-sonner";
import { audioActions, LoadAudioAction, UnloadAudioAction } from "./audio.actions";

export class AudioContext {
    private _worker: Worker | null = null;
    private _pendingRequests = new Map<string, (result: any) => void>();
    private _isWorkerReady = $state(false);


    private _parsingAudio = $state(false);
    private _audioLoaded = $state(false);

    get isWorkerReady() {
        return this._isWorkerReady;
    }

    get parsingAudio() {
        return this._parsingAudio;
    }

    get audioLoaded() {
        return this._audioLoaded;
    }

    async initWorker() {
        if (this._worker) return;

        this._worker = new Worker(
            new URL('$lib/context/audio.worker.ts', import.meta.url),
            { type: 'module' }
        )

        this._worker.onmessage = (event) => {
            const { type, payload, id } = event.data;
            if (type === 'ready') {
                this._isWorkerReady = true;
                return;
            }
            for (const action of audioActions) {
                if (type === action.responseKey) {
                    const resolve = this._pendingRequests.get(id);
                    if (resolve) {
                        resolve(payload);
                        this._pendingRequests.delete(id);
                    }
                    break;
                }
            }
        }
    }

    private sendMessage<T>(type: string, payload?: any): Promise<T> {
        if (!this._worker) {
            throw new Error('Message sent before worker is ready. Call initWorker() first.');
        }
        return new Promise((resolve) => {
            const id = crypto.randomUUID();
            this._pendingRequests.set(id, resolve);
            this._worker?.postMessage({ type, payload: payload ?? undefined, id });
        })
    }

    async loadAudio(file: File) {
        this._parsingAudio = true;

        const arrayBuffer = await file.arrayBuffer();
        this.loadAudioFromArrayBuffer(arrayBuffer);

        this._parsingAudio = false;
        this._audioLoaded = true;
    }

    async loadAudioFromSrc(src: string) {
        this._parsingAudio = true;

        const response = await fetch(src);

        if (!response.ok) {
            toast.error("Failed to download file", { description: "Please try again later" })
            this._parsingAudio = false;
            return;
        }

        await this.loadAudioFromArrayBuffer(await response.arrayBuffer());
    }

    private async loadAudioFromArrayBuffer(arrayBuffer: ArrayBuffer) {
        const uint8Array = new Uint8Array(arrayBuffer);

        const success = await this.sendMessage<boolean>(LoadAudioAction.requestKey, uint8Array);
        if (!success) {
            toast.error("Failed to parse audio", { description: "Are you sure this is a mp3 file?" })
            this._parsingAudio = false;
            return;
        }

        this._parsingAudio = false;
        this._audioLoaded = true;
    }

    async resetAudio() {
        const success = await this.sendMessage<boolean>(UnloadAudioAction.requestKey);
        if (!success) {
            toast.error("Failed to unload audio", { description: "Please try again" })
            return;
        }

        this._audioLoaded = false;
    }
}


const AUDIO_CONTEXT_KEY = Symbol('AUDIO_CONTEXT');

export function setAudioContext() {
    return setContext(AUDIO_CONTEXT_KEY, new AudioContext());
}

export function getAudioContext() {
    return getContext(AUDIO_CONTEXT_KEY) as AudioContext;
}