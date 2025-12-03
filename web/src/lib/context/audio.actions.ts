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

export const LoadAudioAction = AudioAction.create('loadAudio', 'loadAudio');
export const UnloadAudioAction = AudioAction.create('unloadAudio', 'unloadAudio');
export const GetAudioMetadataAction = AudioAction.create('getAudioMetadata', 'getAudioMetadata');

export const audioActions = [
    LoadAudioAction,
    UnloadAudioAction,
    GetAudioMetadataAction,
] as const;