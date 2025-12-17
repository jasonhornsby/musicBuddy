import type { AudioBufferSetup } from "$lib/utils/audioBufferManager";
import type { Visualisation } from "../visualiser.svelte";
import type { BaseVizConfig } from "$lib/types/nodeConfig";
import AudioWorker from './audio.worker.ts?worker';
import type { RenderTree } from "$lib/components/track/pipeline-graph.svelte";
import { SvelteMap } from "svelte/reactivity";

type ParamDef = {
    key: string;
    label: string;
    type: string;
    default: any;
    min?: number;
    max?: number;
    step?: number;
    options?: string[];
};

export class AudioWorkerManager {
    public isReady = $state(false);
    public isAudioLoaded = $state(false);

    private worker: Worker;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    private visualisations: Record<string, Visualisation<any>> = {};
    public renderTree: RenderTree = $state({
        nodes: [],
        edges: []
    });
    public configSchemas: SvelteMap<string, ParamDef[]> = new SvelteMap();

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
                case 'viz_ready':
                    this.visualisations[data.id].setReady(true);
                    break;
                case 'viz_updated':
                    if (!this.visualisations[data.id].ready) {
                        break;
                    }
                    this.visualisations[data.id].draw();
                    break;
                case 'viz_created':
                    this.renderTree = JSON.parse(data.tree) as RenderTree;
                    this.configSchemas.set(data.id, JSON.parse(data.schema) as ParamDef[]);
                    console.log('[TS] Viz created:', { renderTree: this.renderTree, configSchema: this.configSchemas.get(data.id) });
                    break;
                case 'viz_configured':
                    // Config applied, now trigger an update to refresh the viz
                    this.updateVisualizer(data.id);
                    break;
                case 'render_tree_updated':
                    console.log('[TS] Render tree updated:', JSON.parse(data.tree));
                    this.renderTree = JSON.parse(data.tree) as RenderTree;
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

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    public createVisualizer(visualisation: Visualisation<any>) {
        this.worker.postMessage({
            type: 'create_viz',
            id: visualisation.id,
            vizType: visualisation.type,
            buffer: visualisation.uInt8View,
            config: visualisation.config
        });
        this.visualisations[visualisation.id] = visualisation;

        // Wire up config change handler
        visualisation.registerConfigChangeHandler((config) => {
            this.configureVisualizer(visualisation.id, config);
        });
    }

    public configureVisualizer(id: string, config: BaseVizConfig) {
        this.worker.postMessage({
            type: 'configure_viz',
            id: id,
            config: config
        });
    }

    public updateVisualizer(id: string) {
        this.worker.postMessage({
            type: 'update_viz',
            id: id
        });
    }
}