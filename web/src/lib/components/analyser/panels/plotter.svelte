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
      chart1: '#d97341',              // oklch(0.646 0.222 41.116) - orange
      chart1Alpha: 'rgba(217, 115, 65, 0.25)'
    };
    
    // Downsample using min-max method for waveform display
    function downsampleMinMax(float64Array: Float64Array, startIdx: number, endIdx: number, targetPoints: number) {
      if (!float64Array || float64Array.length === 0) return [];
      
      const dataLength = endIdx - startIdx;
      const blockSize = Math.max(1, Math.floor(dataLength / targetPoints));
      const downsampled = [];
      
      for (let i = 0; i < targetPoints; i++) {
        const start = startIdx + (i * blockSize);
        const end = Math.min(start + blockSize, endIdx);
        
        if (start >= endIdx) break;
        
        let min = Infinity;
        let max = -Infinity;
        
        for (let j = start; j < end; j++) {
          if (float64Array[j] < min) min = float64Array[j];
          if (float64Array[j] > max) max = float64Array[j];
        }
        
        // Create two points for min and max to visualize the envelope
        downsampled.push([start, max]);
        downsampled.push([start, min]);
      }
      
      return downsampled;
    }
    
    function getDataForZoom(start: number, end: number) {
      if (!audioContext.decoded) return [];
      
      const totalLength = audioContext.decoded.length;
      const startIdx = Math.floor((start / 100) * totalLength);
      const endIdx = Math.ceil((end / 100) * totalLength);
      const visibleRange = endIdx - startIdx;
      
      // Adjust target points based on zoom level
      // More zoomed in = more detail
      const targetPoints = Math.min(visibleRange, 500);
      
      return downsampleMinMax(audioContext.decoded, startIdx, endIdx, targetPoints);
    }
    
    function initChart() {
      if (!chartContainer || !audioContext.decoded) return;
      
      chart = echarts.init(chartContainer as HTMLElement);

      if (!chart) return;

      
      const option = {
        backgroundColor: colors.background,
        grid: {
          left: 50,
          right: 50,
          top: 30,
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
            const sample = params[0].data[0];
            const amplitude = params[0].data[1].toFixed(4);
            return `Sample: ${sample}<br/>Amplitude: ${amplitude}`;
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
            type: 'line',
            data: getDataForZoom(0, 100),
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
          }
        ]
      };
      
      chart.setOption(option);
      
      // Listen to zoom events and update data resolution
      chart.on('dataZoom', (params) => {
        if (!chart) return;
        const option = chart.getOption() as echarts.EChartsOption;
        const dataZoom = option.dataZoom[0] as any;
        const start = dataZoom.start;
        const end = dataZoom.end;
        
        // Only update if zoom changed significantly (avoid excessive updates)
        if (Math.abs(start - currentZoom.start) > 0.5 || 
            Math.abs(end - currentZoom.end) > 0.5) {
          currentZoom = { start, end };
          
          const newData = getDataForZoom(start, end);
          
          chart.setOption({
            series: [{ data: newData }]
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
