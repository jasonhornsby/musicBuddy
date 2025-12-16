import type { BaseVizConfig, SpectrumConfig } from '$lib/types/nodeConfig';

export abstract class Visualisation<TConfig extends BaseVizConfig = BaseVizConfig> {
    id: string;
    type: string;
    // Width of the output buffer not the visualisation itself
    // TODO: Evaluate if this is the best way to do this
    // We don't want to resise the buffer each time we resize the screen but for big screen width we 
    // still need a large buffer so we have something to display
    // The JS code takes this layout and draws the visualisation. The buffer does not contain the bitmap of the viz
    requestedDatapoints: number;
    ctx: CanvasRenderingContext2D | null = null;

    private buffer: SharedArrayBuffer | null = null;
    private _floatView: Float32Array | null = null;
    private _uInt8View: Uint8Array | null = null;
    private _config: TConfig;
    private onConfigChange?: (config: TConfig) => void;

    get floatView(): Float32Array {
        return this._floatView as Float32Array;
    }

    set floatView(value: Float32Array) {
        this._floatView = value;
    }

    get uInt8View(): Uint8Array {
        return this._uInt8View as Uint8Array;
    }

    set uInt8View(value: Uint8Array) {
        this._uInt8View = value;
    }

    get config(): TConfig {
        return this._config;
    }

    constructor(id: string, type: string, requestedDatapoints: number, defaultConfig: TConfig) {
        this.id = id;
        this.type = type;
        this.requestedDatapoints = requestedDatapoints;
        this._config = defaultConfig;
        this.createBuffer();
    }

    public abstract getBufferSize(): number;
    public abstract draw(): void;

    public getBuffer(): SharedArrayBuffer {
        if (!this.buffer) {
            this.createBuffer();
        }
        return this.buffer as SharedArrayBuffer;
    }

    public createBuffer() {
        this.buffer = new SharedArrayBuffer(this.getBufferSize());
        this.floatView = new Float32Array(this.buffer);
        this.uInt8View = new Uint8Array(this.buffer);
    }

    public registerContext(ctx: CanvasRenderingContext2D) {
        this.ctx = ctx;
    }

    public updateConfig(newConfig: Partial<TConfig>) {
        this._config = { ...this._config, ...newConfig };
        this.onConfigChange?.(this._config);
    }

    public registerConfigChangeHandler(handler: (config: TConfig) => void) {
        this.onConfigChange = handler;
    }
}


import type { WaveformConfig } from '$lib/types/nodeConfig';

export class WaveformVisualisation extends Visualisation<WaveformConfig> {
    constructor(requestedDatapoints: number, config?: Partial<WaveformConfig>) {
        const defaultConfig: WaveformConfig = {
            channel: 'mix',
            ...config
        };
        super(
            `waveform`,
            'waveform',
            requestedDatapoints,
            defaultConfig
        );
    }

    public getBufferSize(): number {
        // 4 bytes per float, 2 floats per sample (min, max)
        return this.requestedDatapoints * 4 * 2;
    }

    public draw(): void {
        if (!this.ctx) {
            throw new Error('Context not registered');
        }

        // Clear the entire canvas before redrawing
        this.ctx.clearRect(0, 0, this.ctx.canvas.width, this.ctx.canvas.height);

        const dpr = window.devicePixelRatio;
        const canvasHeight = this.ctx.canvas.height / dpr;
        const canvasWidth = this.ctx.canvas.width / dpr;

        // Scale configuration
        const scaleWidth = 20;
        const scalePadding = 4;
        const scaleTotalWidth = scaleWidth + scalePadding;
        const tickLength = 6;
        const subTickLength = 3;
        const scaleValues = [1, 0.5, 0, -0.5, -1];
        const subTickValues = [0.75, 0.25, -0.25, -0.75];
        const verticalPadding = 6;

        // Draw scale background (matching track.svelte slate-50/80 background)
        this.ctx.fillStyle = '#f8fafc';
        this.ctx.fillRect(0, 0, scaleTotalWidth, canvasHeight);

        // Draw scale line (matching track.svelte slate-200 border)
        this.ctx.strokeStyle = '#e2e8f0';
        this.ctx.lineWidth = 1;
        this.ctx.beginPath();
        this.ctx.moveTo(scaleTotalWidth, 0);
        this.ctx.lineTo(scaleTotalWidth, canvasHeight);
        this.ctx.stroke();

        // Draw faint horizontal gridlines for main ticks
        this.ctx.strokeStyle = '#f1f5f9';
        this.ctx.lineWidth = 0.5;
        for (const value of scaleValues) {
            const y = ((1 - value) / 2) * (canvasHeight - verticalPadding * 2) + verticalPadding;
            this.ctx.beginPath();
            this.ctx.moveTo(scaleTotalWidth, y);
            this.ctx.lineTo(canvasWidth, y);
            this.ctx.stroke();
        }

        // Draw even fainter horizontal gridlines for sub-ticks
        this.ctx.strokeStyle = '#f8fafc';
        this.ctx.lineWidth = 0.5;
        for (const value of subTickValues) {
            const y = ((1 - value) / 2) * (canvasHeight - verticalPadding * 2) + verticalPadding;
            this.ctx.beginPath();
            this.ctx.moveTo(scaleTotalWidth, y);
            this.ctx.lineTo(canvasWidth, y);
            this.ctx.stroke();
        }

        // Draw tick marks and labels (matching track.svelte slate-800 text)
        this.ctx.fillStyle = '#1e293b';
        this.ctx.font = '6px system-ui, -apple-system, sans-serif';
        this.ctx.textAlign = 'right';

        for (const value of scaleValues) {
            const y = ((1 - value) / 2) * (canvasHeight - verticalPadding * 2) + verticalPadding;

            // Draw tick mark (matching track.svelte slate-200 border)
            this.ctx.strokeStyle = '#cbd5e1';
            this.ctx.beginPath();
            this.ctx.moveTo(scaleTotalWidth - tickLength, y);
            this.ctx.lineTo(scaleTotalWidth, y);
            this.ctx.stroke();
            this.ctx.textBaseline = 'middle';


            // Draw label
            const label = value > 0 ? `+${value}` : value.toString();
            this.ctx.fillText(label, scaleTotalWidth - tickLength - 2, y);
        }

        // Draw sub ticks (no labels) - lighter slate for subtler ticks
        this.ctx.strokeStyle = '#e2e8f0';
        for (const value of subTickValues) {
            const y = ((1 - value) / 2) * (canvasHeight - verticalPadding * 2) + verticalPadding;
            this.ctx.beginPath();
            this.ctx.moveTo(scaleTotalWidth - subTickLength, y);
            this.ctx.lineTo(scaleTotalWidth, y);
            this.ctx.stroke();
        }

        // Draw min/max bars (data is interleaved: [min, max, min, max, ...])
        // Using modern slate-blue color that matches the track design
        this.ctx.strokeStyle = '#475569';
        this.ctx.lineWidth = 1.5;
        this.ctx.lineCap = 'round';

        const numBars = this.floatView.length / 2;
        const waveformWidth = canvasWidth - scaleTotalWidth;
        const barWidth = waveformWidth / numBars;

        let minPoint = Number.POSITIVE_INFINITY;
        let maxPoint = Number.NEGATIVE_INFINITY;

        for (let i = 0; i < numBars; i++) {
            const minVal = this.floatView[i * 2];
            const maxVal = this.floatView[i * 2 + 1];

            // Track smallest and largest values
            if (minVal < minPoint) minPoint = minVal;
            if (maxVal > maxPoint) maxPoint = maxVal;

            // Normalize values (assuming -1 to 1 range) to canvas coordinates
            const yMin = ((1 - minVal) / 2) * canvasHeight;
            const yMax = ((1 - maxVal) / 2) * canvasHeight;

            const x = scaleTotalWidth + i * barWidth + barWidth / 2;

            this.ctx.beginPath();
            this.ctx.moveTo(x, yMin);
            this.ctx.lineTo(x, yMax);
            this.ctx.stroke();
        }

        console.log('Waveform smallest value:', minPoint, 'largest value:', maxPoint);
    }
}

export class SpectrumVisualisation extends Visualisation<SpectrumConfig> {
    constructor(requestedDatapoints: number, config?: Partial<SpectrumConfig>) {
        const defaultConfig: SpectrumConfig = {
            channel: 'mix',
            windowSize: 1024,
            ...config
        };
        super(
            `spectrum`,
            'spectrum',
            requestedDatapoints,
            defaultConfig
        );
    }

    public getBufferSize(): number {
        return this.requestedDatapoints * 4;
    }

    public draw(): void {
        if (!this.ctx) {
            throw new Error('Context not registered');
        }

        // Clear the entire canvas before redrawing
        this.ctx.clearRect(0, 0, this.ctx.canvas.width, this.ctx.canvas.height);

        const dpr = window.devicePixelRatio;
        const canvasHeight = this.ctx.canvas.height / dpr;
        const canvasWidth = this.ctx.canvas.width / dpr;

        // Scale configuration (matching WaveformVisualisation)
        const scaleWidth = 20;
        const scalePadding = 4;
        const scaleTotalWidth = scaleWidth + scalePadding;
        const tickLength = 6;
        const subTickLength = 3;
        const verticalPadding = 6;

        // Draw scale background (matching track.svelte slate-50/80 background)
        this.ctx.fillStyle = '#f8fafc';
        this.ctx.fillRect(0, 0, scaleTotalWidth, canvasHeight);

        // Draw scale line (matching track.svelte slate-200 border)
        this.ctx.strokeStyle = '#e2e8f0';
        this.ctx.lineWidth = 1;
        this.ctx.beginPath();
        this.ctx.moveTo(scaleTotalWidth, 0);
        this.ctx.lineTo(scaleTotalWidth, canvasHeight);
        this.ctx.stroke();

        // Find min/max values for normalization
        let minVal = Number.POSITIVE_INFINITY;
        let maxVal = Number.NEGATIVE_INFINITY;
        for (let i = 0; i < this.floatView.length; i++) {
            const val = this.floatView[i];
            if (val < minVal) minVal = val;
            if (val > maxVal) maxVal = val;
        }

        // Handle edge case where all values are the same
        const range = maxVal - minVal || 1;

        // Draw faint horizontal gridlines
        this.ctx.strokeStyle = '#f1f5f9';
        this.ctx.lineWidth = 0.5;
        const gridSteps = 4;
        for (let i = 0; i <= gridSteps; i++) {
            const y = verticalPadding + (i / gridSteps) * (canvasHeight - verticalPadding * 2);
            this.ctx.beginPath();
            this.ctx.moveTo(scaleTotalWidth, y);
            this.ctx.lineTo(canvasWidth, y);
            this.ctx.stroke();
        }

        // Draw tick marks and labels (matching track.svelte slate-800 text)
        this.ctx.fillStyle = '#1e293b';
        this.ctx.font = '6px system-ui, -apple-system, sans-serif';
        this.ctx.textAlign = 'right';
        this.ctx.textBaseline = 'middle';

        for (let i = 0; i <= gridSteps; i++) {
            const y = verticalPadding + (i / gridSteps) * (canvasHeight - verticalPadding * 2);
            const value = maxVal - (i / gridSteps) * range;

            // Draw tick mark (matching track.svelte slate-200 border)
            this.ctx.strokeStyle = '#cbd5e1';
            this.ctx.beginPath();
            this.ctx.moveTo(scaleTotalWidth - tickLength, y);
            this.ctx.lineTo(scaleTotalWidth, y);
            this.ctx.stroke();

            // Draw label
            const label = value.toFixed(2);
            this.ctx.fillText(label, scaleTotalWidth - tickLength - 2, y);
        }

        // Draw sub ticks (no labels) - lighter slate for subtler ticks
        this.ctx.strokeStyle = '#e2e8f0';
        for (let i = 0; i < gridSteps; i++) {
            const y = verticalPadding + ((i + 0.5) / gridSteps) * (canvasHeight - verticalPadding * 2);
            this.ctx.beginPath();
            this.ctx.moveTo(scaleTotalWidth - subTickLength, y);
            this.ctx.lineTo(scaleTotalWidth, y);
            this.ctx.stroke();
        }

        // Draw the spectrum line plot (using modern slate-blue color)
        this.ctx.strokeStyle = '#475569';
        this.ctx.lineWidth = 1.5;
        this.ctx.lineCap = 'round';
        this.ctx.lineJoin = 'round';

        const numPoints = this.floatView.length;
        const plotWidth = canvasWidth - scaleTotalWidth;
        const pointSpacing = plotWidth / (numPoints - 1 || 1);

        this.ctx.beginPath();
        for (let i = 0; i < numPoints; i++) {
            const val = this.floatView[i];
            const normalizedVal = (val - minVal) / range;
            const x = scaleTotalWidth + i * pointSpacing;
            const y = verticalPadding + (1 - normalizedVal) * (canvasHeight - verticalPadding * 2);

            if (i === 0) {
                this.ctx.moveTo(x, y);
            } else {
                this.ctx.lineTo(x, y);
            }
        }
        this.ctx.stroke();
    }
}