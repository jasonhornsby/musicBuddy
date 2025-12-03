import { getContext, setContext } from "svelte";
import { toast } from "svelte-sonner";

export class AudioContext {


    private _parsingAudio = $state(false);
    private _audioLoaded = $state(false);

    get parsingAudio() {
        return this._parsingAudio;
    }

    get audioLoaded() {
        return this._audioLoaded;
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
        const success = await window.loadAudio(uint8Array);
        if (!success) {
            toast.error("Failed to parse audio", { description: "Are you sure this is a mp3 file?" })
            this._parsingAudio = false;
            return;
        }

        this._parsingAudio = false;
        this._audioLoaded = true;
    }

    resetAudio() {
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