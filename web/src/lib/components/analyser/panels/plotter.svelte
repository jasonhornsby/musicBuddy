<script lang="ts">
    import * as echarts from 'echarts';
	import { getAudioContext } from '$lib/context/audio.svelte';
	import { onMount } from 'svelte';
    
    const audioContext = getAudioContext();
    
    let chartContainer = $state<HTMLElement | undefined>();
    let chart = $state<echarts.ECharts | null>(null);
    let currentZoom = $state({ start: 0, end: 100 });
    
    // Theme colors matching project palette (light mode)
    const colors = {
      background: '#ffffff',
      foreground: '#1a1625',          // oklch(0.129 0.042 264.695) - dark purple
      primary: '#1f1a2e',             // oklch(0.208 0.042 265.755) - primary
      muted: '#f4f4f6',               // oklch(0.968 0.007 247.896) - muted bg
      mutedForeground: '#6b6680',     // oklch(0.554 0.046 257.417)
      border: '#e5e4e9',              // oklch(0.929 0.013 255.508)
      chart1: '#d97341',              // oklch(0.646 0.222 41.116) - orange (left channel)
      chart1Alpha: 'rgba(217, 115, 65, 0.25)',
      chart2: '#4171d9',              // Blue (right channel)
      chart2Alpha: 'rgba(65, 113, 217, 0.25)'
    };
    
    // Downsample using min-max method for waveform display
    // Float64Array contains interleaved stereo data: [L0, R0, L1, R1, L2, R2, ...]
    // startIdx and endIdx are sample indices (not array indices)
    function downsampleMinMax(float64Array: Float64Array, startSample: number, endSample: number, targetPoints: number) {
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
    
    function getDataForZoom(start: number, end: number) {
      if (!audioContext.decoded) return { left: [], right: [] };
      
      // Total number of stereo samples (array length / 2 since interleaved)
      const totalSamples = audioContext.decoded.length / 2;
      const startSample = Math.floor((start / 100) * totalSamples);
      const endSample = Math.ceil((end / 100) * totalSamples);
      const visibleSamples = endSample - startSample;
      
      // Adjust target points based on zoom level
      // More zoomed in = more detail
      const targetPoints = Math.min(visibleSamples, 500);
      
      return downsampleMinMax(audioContext.decoded, startSample, endSample, targetPoints);
    }
    
    function initChart() {
      if (!chartContainer || !audioContext.decoded) return;
      
      chart = echarts.init(chartContainer as HTMLElement);

      if (!chart) return;

      
      const option = {
        backgroundColor: colors.background,
        legend: {
          show: true,
          top: 5,
          textStyle: { color: colors.mutedForeground }
        },
        grid: {
          left: 50,
          right: 50,
          top: 35,
          bottom: 80
        },
        tooltip: {
          trigger: 'axis',
          backgroundColor: colors.background,
          borderColor: colors.border,
          textStyle: {
            color: colors.foreground
          },
          formatter: (params: any) => {
            if (!params || params.length === 0) return '';
            const sample = params[0].data[0];
            let result = `Sample: ${sample}`;
            for (const param of params) {
              const channelName = param.seriesName;
              const amplitude = param.data[1].toFixed(4);
              const color = param.color;
              result += `<br/><span style="color:${color}">●</span> ${channelName}: ${amplitude}`;
            }
            return result;
          }
        },
        xAxis: {
          type: 'value',
          name: 'Sample',
          nameLocation: 'center',
          nameGap: 30,
          nameTextStyle: { color: colors.mutedForeground },
          axisLine: { lineStyle: { color: colors.border } },
          axisLabel: { color: colors.mutedForeground },
          splitLine: { show: false }
        },
        yAxis: {
          type: 'value',
          name: 'Amplitude',
          nameLocation: 'center',
          nameGap: 40,
          nameTextStyle: { color: colors.mutedForeground },
          min: -1,
          max: 1,
          axisLine: { lineStyle: { color: colors.border } },
          axisLabel: { color: colors.mutedForeground },
          splitLine: { lineStyle: { color: colors.muted } }
        },
        dataZoom: [
          {
            type: 'inside',
            start: 0,
            end: 100,
            zoomOnMouseWheel: true,
            moveOnMouseMove: true
          },
          {
            type: 'slider',
            start: 0,
            end: 100,
            height: 30,
            bottom: 10,
            borderColor: colors.border,
            fillerColor: colors.chart1Alpha,
            handleStyle: {
              color: colors.chart1,
              borderColor: colors.chart1
            },
            textStyle: {
              color: colors.mutedForeground
            },
            dataBackground: {
              lineStyle: { color: colors.border },
              areaStyle: { color: colors.muted }
            },
            selectedDataBackground: {
              lineStyle: { color: colors.chart1 },
              areaStyle: { color: colors.chart1Alpha }
            }
          }
        ],
        series: [
          {
            name: 'Left Channel',
            type: 'line',
            data: getDataForZoom(0, 100).left,
            showSymbol: false,
            lineStyle: {
              color: colors.chart1,
              width: 1
            },
            areaStyle: {
              color: colors.chart1Alpha
            },
            sampling: 'lttb',
            large: true,
            largeThreshold: 1000
          },
          {
            name: 'Right Channel',
            type: 'line',
            data: getDataForZoom(0, 100).right,
            showSymbol: false,
            lineStyle: {
              color: colors.chart2,
              width: 1
            },
            areaStyle: {
              color: colors.chart2Alpha
            },
            sampling: 'lttb',
            large: true,
            largeThreshold: 1000
          }
        ]
      };
      
      chart.setOption(option);
      
      // Listen to zoom events and update data resolution
      chart.on('dataZoom', (params) => {
        if (!chart) return;
        const option = chart.getOption() as echarts.EChartsOption;
        const dataZoomArr = option.dataZoom as echarts.DataZoomComponentOption[] | undefined;
        if (!dataZoomArr || dataZoomArr.length === 0) return;
        const dataZoom = dataZoomArr[0] as { start?: number; end?: number };
        const start = dataZoom.start ?? 0;
        const end = dataZoom.end ?? 100;
        
        // Only update if zoom changed significantly (avoid excessive updates)
        if (Math.abs(start - currentZoom.start) > 0.5 || 
            Math.abs(end - currentZoom.end) > 0.5) {
          currentZoom = { start, end };
          
          const newData = getDataForZoom(start, end);
          
          chart.setOption({
            series: [
              { data: newData.left },
              { data: newData.right }
            ]
          });
        }
      });
      
      // Handle window resize
      window.addEventListener('resize', handleResize);
    }
    
    function handleResize() {
      if (chart) {
        chart.resize();
      }
    }

    onMount(() =>{
        initChart();

        return () => {
            if (chart) {
                chart.dispose();
                window.removeEventListener('resize', handleResize);
            }
        };
    });
</script>

<div class="waveform-container" bind:this={chartContainer}></div>

<style>
    .waveform-container {
      width: 100%;
      height: 100%;
      min-height: 0;
      background: var(--background);
      border-radius: var(--radius);
      border: 1px solid var(--border);
    }
</style>
