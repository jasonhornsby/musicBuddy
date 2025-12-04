import { getContext, setContext } from "svelte";
import { toast } from "svelte-sonner";
import { audioActions, GetAudioMetadataAction, GetSpectralFluxAction, LoadAudioAction, LoadParsedAudioAction, UnloadAudioAction, type AudioMetadata } from "./audio.actions";

export class AudioContext {
    private _worker: Worker | null = null;
    private _pendingRequests = new Map<string, (result: any) => void>();
    private _isWorkerReady = $state(false);
    private _useHardwareAcceleration = $state(true);

    private _parsingAudio = $state(false);
    private _audioLoaded = $state(false);
    // Data
    private _decoded = $state<Float64Array | null>(null);
    private _spectralFlux = $state<Float64Array | null>(null);

    // Insights
    private _metadata = $state<AudioMetadata | null>(null);



    get useHardwareAcceleration() {
        return this._useHardwareAcceleration;
    }


    set useHardwareAcceleration(value: boolean) {
        this._useHardwareAcceleration = value;
    }

    get decoded() {
        return this._decoded;
    }

    get isWorkerReady() {
        return this._isWorkerReady;
    }

    get parsingAudio() {
        return this._parsingAudio;
    }

    get audioLoaded() {
        return this._audioLoaded;
    }

    get metadata() {
        return this._metadata;
    }

    get spectralFlux() {
        return this._spectralFlux;
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

    async getAudioMetadata() {
        console.log("Getting audio metadata");
        const result = await this.sendMessage<AudioMetadata>(GetAudioMetadataAction.requestKey);
        console.log("Audio metadata: ", result);
        this._metadata = result;
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

        let result: Float64Array | null = null;
        if (this.useHardwareAcceleration) {
            // Copy the arrayBuffer before decodeAudioData consumes it
            const rawDataCopy = new Uint8Array(arrayBuffer.slice(0));
            try {
                const initialResult = await this.parseAudioJs(arrayBuffer);
                result = initialResult.data;
                const payload = {
                    parsedAudio: new Uint8Array(result),
                    rawData: rawDataCopy,
                }
                // since we parsed this in JS real we need to send the result to the worker
                await this.sendMessage<Float64Array | null>(LoadParsedAudioAction.requestKey, payload);

            } catch (e) {
                toast.error("Failed to parse audio", { description: "Are you sure this is a mp3 file?" })
                this._parsingAudio = false;
                return;
            }

        } else {
            result = await this.sendMessage<Float64Array | null>(LoadAudioAction.requestKey, uint8Array);
            if (!result) {
                toast.error("Failed to parse audio", { description: "Are you sure this is a mp3 file?" })
                this._parsingAudio = false;
                return;
            }
        }

        this._decoded = result;

        this._parsingAudio = false;
        this._audioLoaded = true;
    }

    private async parseAudioJs(arrayBuffer: ArrayBuffer) {
        const ctx = new (window.AudioContext || window.webkitAudioContext)();
        // Using the AudioContext compared to the work is insanely faster. HWA FTW!
        const decoded = await ctx.decodeAudioData(arrayBuffer)

        decoded.sampleRate

        // Need to convert the decoded audio to a Float64Array
        // The channels need to be interleaved
        const interleaved = new Float64Array(decoded.length * 2);
        for (let i = 0; i < decoded.length; i++) {
            interleaved[i * 2] = decoded.getChannelData(0)[i];
            interleaved[i * 2 + 1] = decoded.getChannelData(1)[i];
        }
        return {
            data: interleaved,
            meta: { sampleRate: decoded.sampleRate, channels: decoded.numberOfChannels }
        };
    }

    async resetAudio() {
        const success = await this.sendMessage<boolean>(UnloadAudioAction.requestKey);
        if (!success) {
            toast.error("Failed to unload audio", { description: "Please try again" })
            return;
        }

        this._decoded = null;
        this._audioLoaded = false;
    }

    async getSpectralFlux() {
        if (this._spectralFlux) return;
        const result = await this.sendMessage<Float64Array | null>(GetSpectralFluxAction.requestKey);
        if (!result) {
            toast.error("Failed to get spectral flux", { description: "Please try again" })
            return;
        }
        this._spectralFlux = result;
    }
}


const AUDIO_CONTEXT_KEY = Symbol('AUDIO_CONTEXT');

export function setAudioContext() {
    return setContext(AUDIO_CONTEXT_KEY, new AudioContext());
}

export function getAudioContext() {
    return getContext(AUDIO_CONTEXT_KEY) as AudioContext;
}