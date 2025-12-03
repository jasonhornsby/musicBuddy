import * as echarts from 'echarts';

// Shared theme colors matching project palette (light mode)
export const chartColors = {
  background: '#ffffff',
  foreground: '#1a1625',          // oklch(0.129 0.042 264.695) - dark purple
  primary: '#1f1a2e',             // oklch(0.208 0.042 265.755) - primary
  muted: '#f4f4f6',               // oklch(0.968 0.007 247.896) - muted bg
  mutedForeground: '#6b6680',     // oklch(0.554 0.046 257.417)
  border: '#e5e4e9',              // oklch(0.929 0.013 255.508)
} as const;

// Chart-specific color palettes
export const chartPalettes = {
  waveform: {
    left: '#d97341',              // oklch(0.646 0.222 41.116) - orange
    leftAlpha: 'rgba(217, 115, 65, 0.25)',
    right: '#4171d9',             // Blue
    rightAlpha: 'rgba(65, 113, 217, 0.25)'
  },
  spectralFlux: {
    flux: '#2dd4bf',              // Teal
    fluxAlpha: 'rgba(45, 212, 191, 0.25)',
    threshold: '#a855f7',         // Purple accent
    thresholdAlpha: 'rgba(168, 85, 247, 0.15)'
  }
} as const;

// Common tooltip configuration
export function createTooltipConfig() {
  return {
    trigger: 'axis' as const,
    backgroundColor: chartColors.background,
    borderColor: chartColors.border,
    textStyle: {
      color: chartColors.foreground
    }
  };
}

// Common grid configuration
export function createGridConfig(options?: { leftGap?: number; rightGap?: number }) {
  return {
    left: options?.leftGap ?? 50,
    right: options?.rightGap ?? 50,
    top: 35,
    bottom: 80
  };
}

// Common legend configuration
export function createLegendConfig() {
  return {
    show: true,
    top: 5,
    textStyle: { color: chartColors.mutedForeground }
  };
}

// Common x-axis configuration
export function createXAxisConfig(options: { 
  name: string; 
  max: number;
  nameGap?: number;
}) {
  return {
    type: 'value' as const,
    name: options.name,
    nameLocation: 'center' as const,
    nameGap: options.nameGap ?? 30,
    nameTextStyle: { color: chartColors.mutedForeground },
    axisLine: { lineStyle: { color: chartColors.border } },
    axisLabel: { color: chartColors.mutedForeground },
    splitLine: { show: false },
    min: 0,
    max: options.max
  };
}

// Common y-axis configuration
export function createYAxisConfig(options: {
  name: string;
  min: number;
  max: number;
  nameGap?: number;
  labelFormatter?: (value: number) => string;
}) {
  return {
    type: 'value' as const,
    name: options.name,
    nameLocation: 'center' as const,
    nameGap: options.nameGap ?? 40,
    nameTextStyle: { color: chartColors.mutedForeground },
    min: options.min,
    max: options.max,
    axisLine: { lineStyle: { color: chartColors.border } },
    axisLabel: { 
      color: chartColors.mutedForeground,
      ...(options.labelFormatter && { formatter: options.labelFormatter })
    },
    splitLine: { lineStyle: { color: chartColors.muted } }
  };
}

// Common dataZoom configuration
export function createDataZoomConfig(accentColor: string, accentColorAlpha: string) {
  return [
    {
      type: 'inside' as const,
      start: 0,
      end: 100,
      zoomOnMouseWheel: true,
      moveOnMouseMove: true
    },
    {
      type: 'slider' as const,
      start: 0,
      end: 100,
      height: 30,
      bottom: 10,
      borderColor: chartColors.border,
      fillerColor: accentColorAlpha,
      handleStyle: {
        color: accentColor,
        borderColor: accentColor
      },
      textStyle: {
        color: chartColors.mutedForeground
      },
      dataBackground: {
        lineStyle: { color: chartColors.border },
        areaStyle: { color: chartColors.muted }
      },
      selectedDataBackground: {
        lineStyle: { color: accentColor },
        areaStyle: { color: accentColorAlpha }
      }
    }
  ];
}

// Common series configuration for line charts
export function createLineSeriesConfig(options: {
  name: string;
  data: [number, number][];
  color: string;
  colorAlpha: string;
  lineWidth?: number;
  withGradient?: boolean;
}) {
  const baseConfig = {
    name: options.name,
    type: 'line' as const,
    data: options.data,
    showSymbol: false,
    lineStyle: {
      color: options.color,
      width: options.lineWidth ?? 1
    },
    sampling: 'lttb' as const,
    large: true,
    largeThreshold: 1000
  };

  if (options.withGradient) {
    return {
      ...baseConfig,
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: options.colorAlpha },
          { offset: 1, color: options.colorAlpha.replace(/[\d.]+\)$/, '0.02)') }
        ])
      }
    };
  }

  return {
    ...baseConfig,
    areaStyle: {
      color: options.colorAlpha
    }
  };
}

// Zoom change detection helper
export function hasZoomChanged(
  currentZoom: { start: number; end: number },
  newStart: number,
  newEnd: number,
  threshold = 0.5
): boolean {
  return Math.abs(newStart - currentZoom.start) > threshold || 
         Math.abs(newEnd - currentZoom.end) > threshold;
}

// Extract zoom values from chart option
export function getZoomFromOption(chart: echarts.ECharts): { start: number; end: number } {
  const option = chart.getOption() as echarts.EChartsOption;
  const dataZoomArr = option.dataZoom as echarts.DataZoomComponentOption[] | undefined;
  if (!dataZoomArr || dataZoomArr.length === 0) {
    return { start: 0, end: 100 };
  }
  const dataZoom = dataZoomArr[0] as { start?: number; end?: number };
  return {
    start: dataZoom.start ?? 0,
    end: dataZoom.end ?? 100
  };
}

// Create resize handler for chart
export function createResizeHandler(chartGetter: () => echarts.ECharts | null) {
  return () => {
    const chart = chartGetter();
    if (chart) {
      chart.resize();
    }
  };
}

// Downsample mono audio data using min-max method
// Returns array of [index, value] points for chart visualization
export function downsampleMinMaxMono(
  float64Array: Float64Array,
  startIdx: number,
  endIdx: number,
  targetPoints: number
): [number, number][] {
  if (!float64Array || float64Array.length === 0) return [];
  
  const numSamples = endIdx - startIdx;
  const blockSize = Math.max(1, Math.floor(numSamples / targetPoints));
  const result: [number, number][] = [];
  
  for (let i = 0; i < targetPoints; i++) {
    const blockStart = startIdx + (i * blockSize);
    const blockEnd = Math.min(blockStart + blockSize, endIdx);
    
    if (blockStart >= endIdx) break;
    
    let min = Infinity;
    let max = -Infinity;
    
    for (let idx = blockStart; idx < blockEnd; idx++) {
      const val = float64Array[idx];
      if (val < min) min = val;
      if (val > max) max = val;
    }
    
    // Create two points for min and max to visualize the envelope
    result.push([blockStart, max]);
    result.push([blockStart, min]);
  }
  
  return result;
}

// Downsample interleaved stereo audio data using min-max method
// Float64Array contains interleaved stereo data: [L0, R0, L1, R1, L2, R2, ...]
// startSample and endSample are sample indices (not array indices)
// Returns { left, right } arrays of [index, value] points for chart visualization
export function downsampleMinMaxStereo(
  float64Array: Float64Array,
  startSample: number,
  endSample: number,
  targetPoints: number
): { left: [number, number][]; right: [number, number][] } {
  if (!float64Array || float64Array.length === 0) return { left: [], right: [] };
  
  const numSamples = endSample - startSample;
  const blockSize = Math.max(1, Math.floor(numSamples / targetPoints));
  const leftChannel: [number, number][] = [];
  const rightChannel: [number, number][] = [];
  
  for (let i = 0; i < targetPoints; i++) {
    const blockStart = startSample + (i * blockSize);
    const blockEnd = Math.min(blockStart + blockSize, endSample);
    
    if (blockStart >= endSample) break;
    
    let leftMin = Infinity;
    let leftMax = -Infinity;
    let rightMin = Infinity;
    let rightMax = -Infinity;
    
    for (let sampleIdx = blockStart; sampleIdx < blockEnd; sampleIdx++) {
      // Each sample has 2 values: left at even index, right at odd index
      const arrayIdx = sampleIdx * 2;
      const leftVal = float64Array[arrayIdx];
      const rightVal = float64Array[arrayIdx + 1];
      
      if (leftVal < leftMin) leftMin = leftVal;
      if (leftVal > leftMax) leftMax = leftVal;
      if (rightVal < rightMin) rightMin = rightVal;
      if (rightVal > rightMax) rightMax = rightVal;
    }
    
    // Create two points for min and max to visualize the envelope
    leftChannel.push([blockStart, leftMax]);
    leftChannel.push([blockStart, leftMin]);
    rightChannel.push([blockStart, rightMax]);
    rightChannel.push([blockStart, rightMin]);
  }
  
  return { left: leftChannel, right: rightChannel };
}
