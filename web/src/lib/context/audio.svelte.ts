import { AudioBufferManager } from "$lib/utils/audioBufferManager";
import { AudioWorkerManager } from "$lib/utils/audioWorker/audio-worker-manager.svelte";
import { getContext, setContext } from "svelte";


export class AudioContext {
    public isParsingAudio = $state(false);

    private workerManager = new AudioWorkerManager();
    private bufferManager = new AudioBufferManager();

    public get isWorkerReady() {
        return this.workerManager.isReady;
    }

    public get isAudioLoaded() {
        return this.workerManager.isAudioLoaded;
    }

    public async loadAudioFile(audioFile: File) {
        if (!this.isWorkerReady) {
            throw new Error('Worker not ready yet');
        }
        this.isParsingAudio = true;
        const bufferSetup = await this.bufferManager.loadAudioFile(audioFile);
        this.workerManager.sendAudioData(bufferSetup);
        this.isParsingAudio = false;
    }

    public async loadAudioFromSrc(src: string) {
        if (!this.isWorkerReady) {
            throw new Error('Worker not ready yet');
        }
        this.isParsingAudio = true;
        const bufferSetup = await this.bufferManager.loadAudioFromSrc(src);
        this.workerManager.sendAudioData(bufferSetup);
        this.isParsingAudio = false;
    }
}


const AUDIO_CONTEXT_KEY = Symbol('AUDIO_CONTEXT');

export function setAudioContext() {
    return setContext(AUDIO_CONTEXT_KEY, new AudioContext());
}

export function getAudioContext() {
    return getContext(AUDIO_CONTEXT_KEY) as AudioContext;
}