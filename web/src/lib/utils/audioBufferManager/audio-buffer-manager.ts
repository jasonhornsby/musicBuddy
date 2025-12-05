export interface AudioBufferSetup {
    /**
     * The raw MP3 data as a SharedArrayBuffer
     * Used to extract meta information from the mp3 file
     */
    rawMp3SAB: SharedArrayBuffer;
    rawMp3Size: number;


    decodedChannelSABs: SharedArrayBuffer[];

    numChannels: number;
    numSamples: number;
    sampleRate: number;
    duration: number;
}

export class AudioLoadError extends Error {
    constructor() {
        super("Failed to load audio from url");
        this.name = 'AudioLoadError';
    }
}

export class AudioDecodeError extends Error {
    constructor() {
        super("Failed to decode audio");
        this.name = 'AudioDecodeError';
    }
}


export class AudioBufferManager {
    private audioContext: AudioContext;

    constructor() {
        this.audioContext = new AudioContext();
    }

    async loadAudioFile(audioFile: File) {
        const arrayBuffer = await audioFile.arrayBuffer();
        return this.loadAudio(arrayBuffer);
    }

    async loadAudioFromSrc(src: string) {
        const response = await fetch(src);
        if (!response.ok) {
            throw new AudioLoadError();
        }
        const arrayBuffer = await response.arrayBuffer();
        return this.loadAudio(arrayBuffer);
    }

    async loadAudio(arrayBuffer: ArrayBuffer): Promise<AudioBufferSetup> {
        const rawMp3ArrayBuffer = arrayBuffer;
        const rawMp3Size = rawMp3ArrayBuffer.byteLength;

        // Setup SAB for raw Mp3 data
        const rawMp3SAB = new SharedArrayBuffer(rawMp3Size);
        const rawMp3View = new Uint8Array(rawMp3SAB);
        rawMp3View.set(new Uint8Array(rawMp3ArrayBuffer));

        // Decode the audio, need to copy to avoid array being consumed by decodeAudioData
        let audioBuffer: AudioBuffer;
        try {
            audioBuffer = await this.audioContext.decodeAudioData(rawMp3ArrayBuffer.slice(0));
        } catch (error) {
            throw new AudioDecodeError();
        }

        const numChannels = audioBuffer.numberOfChannels;
        const numSamples = audioBuffer.length;
        const sampleRate = audioBuffer.sampleRate;
        const duration = numSamples / sampleRate;


        console.log(`Decoded: ${numChannels} channels, ${numSamples} samples @ ${sampleRate}Hz`);

        const decodedChannelSABs: SharedArrayBuffer[] = [];

        for (let channel = 0; channel < numChannels; channel++) {
            const channelData = audioBuffer.getChannelData(channel);
            // Create SAB for channel data
            const channelSAB = new SharedArrayBuffer(numSamples * 4);
            const channelView = new Float32Array(channelSAB);

            channelView.set(channelData);
            decodedChannelSABs.push(channelSAB);

            console.log(`Decoded channel ${channel}: ${channelSAB.byteLength} bytes`);
        }

        return {
            rawMp3SAB,
            rawMp3Size,
            decodedChannelSABs,
            numChannels,
            numSamples,
            sampleRate,
            duration,
        }
    }

    public createViews(bufferSetup: AudioBufferSetup) {
        return bufferSetup.decodedChannelSABs.map((sab) => new Float32Array(sab));
    }

}