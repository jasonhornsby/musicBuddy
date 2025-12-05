import { AudioBufferManager, AudioBufferView, type AudioBufferSetup } from "$lib/utils/audioBufferManager";
import { AudioWorkerManager } from "$lib/utils/audioWorker/audio-worker-manager.svelte";
import { getContext, setContext } from "svelte";

export interface AudioInfo {
    numChannels: number;
    numSamples: number;
    sampleRate: number;
    duration: number;
}

export class AudioContext {
    // Lifecycle states
    public isParsingAudio = $state(false);

    // Audio data
    private channelViews = $state<Float32Array[] | null>(null);
    private audioInfo = $state<AudioInfo | null>(null);

    // Managers
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
        this.setAudioData(bufferSetup);
    }

    public async loadAudioFromSrc(src: string) {
        if (!this.isWorkerReady) {
            throw new Error('Worker not ready yet');
        }
        this.isParsingAudio = true;
        const bufferSetup = await this.bufferManager.loadAudioFromSrc(src);
        this.workerManager.sendAudioData(bufferSetup);
        this.setAudioData(bufferSetup);
    }

    private setAudioData(bufferSetup: AudioBufferSetup) {
        this.audioInfo = {
            numChannels: bufferSetup.numChannels,
            numSamples: bufferSetup.numSamples,
            sampleRate: bufferSetup.sampleRate,
            duration: bufferSetup.duration,
        }
        this.channelViews = this.bufferManager.createViews(bufferSetup);
        this.isParsingAudio = false;
    }

    public getAudioBufferView() {
        if (!this.channelViews || !this.audioInfo) {
            throw new Error("Audio data not loaded");
        }
        return new AudioBufferView(this.channelViews, this.audioInfo);
    }
}


const AUDIO_CONTEXT_KEY = Symbol('AUDIO_CONTEXT');

export function setAudioContext() {
    return setContext(AUDIO_CONTEXT_KEY, new AudioContext());
}

export function getAudioContext() {
    return getContext(AUDIO_CONTEXT_KEY) as AudioContext;
}