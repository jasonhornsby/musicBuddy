import type { AudioBufferSetup } from "$lib/utils/audioBufferManager";
import AudioWorker from './audio.worker.ts?worker';

export class AudioWorkerManager {
    private worker: Worker;

    constructor() {
        this.worker = new AudioWorker();

        this.worker.onmessage = (event: MessageEvent) => {
            const { type, ...data } = event.data;

            console.log('Message received from worker:', { type, data });

            switch (type) {
                case 'worker_ready':
                    this.onWorkerReady();
                    break;
                case 'worker_error':
                    break
                case 'audio_loaded':
                    break;
                default:
                    console.warn(`Unknown message type: ${type}`);
            }
        }

        this.worker.onerror = (event: ErrorEvent) => {
            console.error(event);
            console.error('Worker error:', event);
        }
    }

    private onWorkerReady() {
        console.log('Worker ready');
    }

    public sendAudioData(bufferSetup: AudioBufferSetup) {
        this.worker.postMessage({
            type: 'load_audio',
            rawMP3Buffer: bufferSetup.rawMp3SAB,
            rawMP3Size: bufferSetup.rawMp3Size,
            decodedBuffers: bufferSetup.decodedChannelSABs,
            numChannels: bufferSetup.numChannels,
            numSamples: bufferSetup.numSamples,
            sampleRate: bufferSetup.sampleRate,
            duration: bufferSetup.duration
        })
    }

    public terminate() {
        this.worker.terminate();
    }
}