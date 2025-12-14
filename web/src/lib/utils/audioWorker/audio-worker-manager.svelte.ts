import type { AudioBufferSetup } from "$lib/utils/audioBufferManager";
import type { Visualisation } from "../visualiser";
import AudioWorker from './audio.worker.ts?worker';

export class AudioWorkerManager {
    public isReady = $state(false);
    public isAudioLoaded = $state(false);

    private worker: Worker;
    private visualisations: Record<string, Visualisation> = {};

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
                case 'viz_updated':
                    this.visualisations[data.id].draw();
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

    public createVisualizer(visualisation: Visualisation) {
        this.worker.postMessage({
            type: 'create_viz',
            id: visualisation.id,
            vizType: visualisation.type,
            buffer: visualisation.uInt8View
        });
        this.visualisations[visualisation.id] = visualisation;
    }

    public updateVisualizer(id: string) {
        this.worker.postMessage({
            type: 'update_viz',
            id: id
        });
    }
}