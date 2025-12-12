import type { AudioBufferSetup } from "$lib/utils/audioBufferManager";
import AudioWorker from './audio.worker.ts?worker';

export class AudioWorkerManager {
    public isReady = $state(false);
    public isAudioLoaded = $state(false);

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
                    this.isAudioLoaded = true;
                    console.log('[TS] Audio loaded:', data);
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
        this.isReady = true;
    }

    public sendAudioData(bufferSetup: AudioBufferSetup) {
        const decodedViews = bufferSetup.decodedChannelSABs.map((sab) => new Float32Array(sab));
        this.worker.postMessage({
            type: 'load_audio',
            rawMP3Buffer: bufferSetup.rawMp3SAB,
            rawMP3Size: bufferSetup.rawMp3Size,
            decodedBuffers: decodedViews,
            numChannels: bufferSetup.numChannels,
            numSamples: bufferSetup.numSamples,
            sampleRate: bufferSetup.sampleRate,
            duration: bufferSetup.duration
        })
    }

    public terminate() {
        this.worker.terminate();
    }

    public createVisualizer(id: string, type: 'waveform', width: number) {
        // Create SAB for the output buffer
        // What is the layout here?
        // If we want to draw min - max bars we need to double the size of the bufer
        // Buf if the zoom level is to high we need to draw the points directly, halving the bufer size
        const sab = new SharedArrayBuffer(width * 4 * 2);

        // Create JS views into the SAB
        const floatView = new Float32Array(sab);
        const uInt8View = new Uint8Array(sab);

        // Create visualizer
        this.worker.postMessage({
            type: 'create_viz',
            id: id,
            vizType: type,
            buffer: uInt8View
        });

        return floatView;
    }

    public updateVisualizer(id: string) {
        this.worker.postMessage({
            type: 'update_viz',
            id: id
        });
    }
}