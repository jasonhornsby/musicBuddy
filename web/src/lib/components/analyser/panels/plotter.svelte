<script lang="ts">
    import * as echarts from 'echarts';
	import { getAudioContext } from '$lib/context/audio.svelte';
	import { onMount } from 'svelte';
    
    const audioContext = getAudioContext();
    
    let chartContainer = $state();
    let chart = $state<echarts.ECharts | null>(null);
    let currentZoom = $state({ start: 0, end: 100 });
    
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
      console.log('start, end', start, end)
      
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
        backgroundColor: '#1a1a1a',
        grid: {
          left: 50,
          right: 50,
          top: 30,
          bottom: 80
        },
        tooltip: {
          trigger: 'axis',
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
          axisLine: { lineStyle: { color: '#4b5563' } },
          splitLine: { show: false }
        },
        yAxis: {
          type: 'value',
          name: 'Amplitude',
          nameLocation: 'center',
          nameGap: 40,
          min: -1,
          max: 1,
          axisLine: { lineStyle: { color: '#4b5563' } },
          splitLine: { lineStyle: { color: '#374151' } }
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
            borderColor: '#4b5563',
            fillerColor: 'rgba(59, 130, 246, 0.2)',
            handleStyle: {
              color: '#3b82f6'
            },
            textStyle: {
              color: '#9ca3af'
            }
          }
        ],
        series: [
          {
            type: 'line',
            data: getDataForZoom(0, 100),
            showSymbol: false,
            lineStyle: {
              color: '#3b82f6',
              width: 1
            },
            areaStyle: {
              color: 'rgba(59, 130, 246, 0.3)'
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
      height: 400px;
      background: #1a1a1a;
      border-radius: 8px;
    }
</style>
