class AudioAction<T, R> {
    requestKey: string;
    responseKey: string;
    actionKey: string;

    constructor(private key: string, actionKey: string) {
        this.requestKey = `${key}Request`;
        this.responseKey = `${key}Response`;
        this.actionKey = actionKey;
    }

    static create<T, R>(key: string, actionKey: string): AudioAction<T, R> {
        return new AudioAction(key, actionKey);
    }
}

export const LoadAudioAction = AudioAction.create<Float64Array, void>('loadAudio', 'loadAudio');
export const UnloadAudioAction = AudioAction.create('unloadAudio', 'unloadAudio');

export type AudioMetadata = {
    sampleRate: number;
    channels: number;
    durationMs: number;
    decodedBitrate: number;
    metadata: {
        name: string;
        artist: string;
        album: string;
        year: number;
        format: string;
    };
}
export const GetAudioMetadataAction = AudioAction.create<void, AudioMetadata>('getAudioMetadata', 'getAudioMetadata');

export const GetSpectralFluxAction = AudioAction.create<void, Float64Array>('getSpectralFlux', 'getSpectralFlux');

export const audioActions = [
    LoadAudioAction,
    UnloadAudioAction,
    GetAudioMetadataAction,
    GetSpectralFluxAction,
] as const;