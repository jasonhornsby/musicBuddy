interface DownsampleOptions {
    targetPoints: number;
    slice: { start: number; end: number } | null;
}

interface WaveformData {
    min: Float32Array;
    max: Float32Array;
    xAxis: Float32Array;
}

export function downsampleWaveform(
    data: Float32Array,
    options: DownsampleOptions
): WaveformData {
    const { targetPoints, slice } = options;
    const totalLength = data.length;

    // 1. Determine the Start and End indices based on the slice option
    // If slice is null, use the full array.
    // We also clamp values to ensure we don't access out of bounds.
    const startIndex = slice ? Math.max(0, slice.start) : 0;
    const endIndex = slice ? Math.min(totalLength, slice.end) : totalLength;
    const effectiveLength = endIndex - startIndex;

    // 2. Handle Invalid or Empty Range
    if (effectiveLength <= 0) {
        return {
            min: new Float32Array(0),
            max: new Float32Array(0),
            xAxis: new Float32Array(0)
        };
    }

    // 3. Handle Case where the Slice is smaller than Target (No downsampling needed)
    // We just return the raw data within the slice range.
    if (effectiveLength <= targetPoints) {
        const xAxis = new Float32Array(effectiveLength);
        const minOutput = new Float32Array(effectiveLength);
        const maxOutput = new Float32Array(effectiveLength);

        for (let i = 0; i < effectiveLength; i++) {
            const originalIndex = startIndex + i;
            xAxis[i] = originalIndex;
            minOutput[i] = data[originalIndex];
            maxOutput[i] = data[originalIndex];
        }

        return {
            min: minOutput,
            max: maxOutput,
            xAxis
        };
    }

    // 4. Downsampling Logic
    const minOutput = new Float32Array(targetPoints);
    const maxOutput = new Float32Array(targetPoints);
    const xAxis = new Float32Array(targetPoints);

    // The step is now based on the length of the SLICE, not the whole file
    const step = effectiveLength / targetPoints;

    for (let i = 0; i < targetPoints; i++) {
        // Calculate relative bucket positions
        const relativeStart = Math.floor(i * step);
        const relativeEnd = Math.floor((i + 1) * step);

        // Convert to absolute positions in the source data
        const start = startIndex + relativeStart;
        const end = startIndex + relativeEnd;

        // Optimization: Calculate the loop limit once per bucket
        // Ensure we don't read past the defined end of the slice
        const scanEnd = (end > endIndex) ? endIndex : end;

        // Center the X point in the middle of the bucket (Absolute Sample Index)
        xAxis[i] = Math.floor((start + scanEnd - 1) / 2);

        // Initialize with inverted Infinity
        let currentMin = Infinity;
        let currentMax = -Infinity;

        // Inner Loop: Scan the bucket
        for (let j = start; j < scanEnd; j++) {
            const value = data[j];
            if (value < currentMin) currentMin = value;
            if (value > currentMax) currentMax = value;
        }

        // Fallback: If a bucket was somehow empty (e.g. step < 1 due to rounding)
        if (currentMin === Infinity) {
            // If the bucket is empty, we just grab the single point at start
            // or default to 0 if totally out of bounds (shouldn't happen)
            const fallback = data[start] || 0;
            currentMin = fallback;
            currentMax = fallback;
        }

        minOutput[i] = currentMin;
        maxOutput[i] = currentMax;
    }

    return { min: minOutput, max: maxOutput, xAxis };
}