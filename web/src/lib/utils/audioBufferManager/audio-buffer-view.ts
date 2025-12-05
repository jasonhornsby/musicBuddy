import type { AudioInfo } from "$lib/context/audio.svelte";
import { downsampleWaveform } from "../timeseries/downsample";

export class AudioBufferView {
    private channelViews: Float32Array[];
    private audioInfo: AudioInfo;

    constructor(channelViews: Float32Array[], audioInfo: AudioInfo) {
        this.channelViews = channelViews;
        this.audioInfo = audioInfo;
    }

    public get numSamples() {
        return this.audioInfo.numSamples;
    }

    public getChannelViews() {
        return this.channelViews;
    }

    public getAudioInfo() {
        return this.audioInfo;
    }

    public getDownsampledMinMax(channel: number, targetPoints: number, slice: { start: number, end: number } | null = null): [Float32Array, Float32Array, Float32Array] {
        const channelView = this.channelViews[channel];
        const minMax = downsampleWaveform(channelView, { targetPoints, slice });

        return [minMax.xAxis, minMax.min, minMax.max];
    }
}