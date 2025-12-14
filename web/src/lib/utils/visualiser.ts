export abstract class Visualisation {
    id: string;
    type: string;
    // Width of the output buffer not the visualisation itself
    // TODO: Evaluate if this is the best way to do this
    // We don't want to resise the buffer each time we resize the screen but for big screen width we 
    // still need a large buffer so we have something to display
    // The JS code takes this layout and draws the visualisation. The buffer does not contain the bitmap of the viz
    width: number;
    ctx: CanvasRenderingContext2D | null = null;

    private buffer: SharedArrayBuffer | null = null;
    private _floatView: Float32Array | null = null;
    private _uInt8View: Uint8Array | null = null;

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

    constructor(id: string, type: string, width: number) {
        this.id = id;
        this.type = type;
        this.width = width;
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
}


export class WaveformVisualisation extends Visualisation {
    constructor(width: number) {
        super(`waveform-${Math.random().toString(36).substring(2, 15)}`, 'waveform', width);
    }

    public getBufferSize(): number {
        return this.width * 4 * 2;
    }

    public draw(): void {
        if (!this.ctx) {
            throw new Error('Context not registered');
        }

        const dpr = window.devicePixelRatio;
        const canvasHeight = this.ctx.canvas.height / dpr;
        const canvasWidth = this.ctx.canvas.width / dpr;

        // Scale configuration
        const scaleWidth = 40;
        const tickLength = 6;
        const scaleValues = [1, 0.5, 0, -0.5, -1];

        // Draw scale background
        this.ctx.fillStyle = '#1a1a1a';
        this.ctx.fillRect(0, 0, scaleWidth, canvasHeight);

        // Draw scale line
        this.ctx.strokeStyle = '#666';
        this.ctx.lineWidth = 1;
        this.ctx.beginPath();
        this.ctx.moveTo(scaleWidth, 0);
        this.ctx.lineTo(scaleWidth, canvasHeight);
        this.ctx.stroke();

        // Draw tick marks and labels
        this.ctx.fillStyle = '#aaa';
        this.ctx.font = '10px monospace';
        this.ctx.textAlign = 'right';

        for (const value of scaleValues) {
            const y = ((1 - value) / 2) * canvasHeight;

            // Draw tick mark
            this.ctx.beginPath();
            this.ctx.moveTo(scaleWidth - tickLength, y);
            this.ctx.lineTo(scaleWidth, y);
            this.ctx.stroke();

            // Adjust text baseline so edge labels aren't cut off
            if (value === 1) {
                this.ctx.textBaseline = 'top';
            } else if (value === -1) {
                this.ctx.textBaseline = 'bottom';
            } else {
                this.ctx.textBaseline = 'middle';
            }

            // Draw label
            const label = value > 0 ? `+${value}` : value.toString();
            this.ctx.fillText(label, scaleWidth - tickLength - 2, y);
        }

        // Draw min/max bars (data is interleaved: [min, max, min, max, ...])
        this.ctx.strokeStyle = '#00ff88';
        this.ctx.lineWidth = 1;

        const numBars = this.floatView.length / 2;
        const waveformWidth = canvasWidth - scaleWidth;
        const barWidth = waveformWidth / numBars;

        for (let i = 0; i < numBars; i++) {
            const minVal = this.floatView[i * 2];
            const maxVal = this.floatView[i * 2 + 1];

            // Normalize values (assuming -1 to 1 range) to canvas coordinates
            const yMin = ((1 - minVal) / 2) * canvasHeight;
            const yMax = ((1 - maxVal) / 2) * canvasHeight;

            const x = scaleWidth + i * barWidth + barWidth / 2;

            this.ctx.beginPath();
            this.ctx.moveTo(x, yMin);
            this.ctx.lineTo(x, yMax);
            this.ctx.stroke();
        }
    }
}