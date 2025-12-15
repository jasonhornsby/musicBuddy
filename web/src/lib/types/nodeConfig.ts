export type ChannelMode = 'left' | 'right' | 'mix';

export interface BaseVizConfig {
    channel: ChannelMode;
}

export interface WaveformConfig extends BaseVizConfig { }
