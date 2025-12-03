<script lang="ts">
    import * as echarts from 'echarts';
    import { getAudioContext } from '$lib/context/audio.svelte';
    import { onMount } from 'svelte';
    import {
      chartColors,
      chartPalettes,
      createTooltipConfig,
      createGridConfig,
      createLegendConfig,
      createXAxisConfig,
      createYAxisConfig,
      createDataZoomConfig,
      createLineSeriesConfig,
      hasZoomChanged,
      getZoomFromOption,
      createResizeHandler,
      downsampleMinMaxStereo
    } from '$lib/utils/chartConfig';
    
    const audioContext = getAudioContext();
    const palette = chartPalettes.waveform;
    
    let chartContainer = $state<HTMLElement | undefined>();
    let chart = $state<echarts.ECharts | null>(null);
    let currentZoom = $state({ start: 0, end: 100 });
    
    function getDataForZoom(start: number, end: number) {
      if (!audioContext.decoded) return { left: [], right: [] };
      
      // Total number of stereo samples (array length / 2 since interleaved)
      const totalSamples = audioContext.decoded.length / 2;
      const startSample = Math.floor((start / 100) * totalSamples);
      const endSample = Math.ceil((end / 100) * totalSamples);
      const visibleSamples = endSample - startSample;
      
      // Adjust target points based on zoom level
      // More zoomed in = more detail
      const targetPoints = Math.min(visibleSamples, 700);
      
      return downsampleMinMaxStereo(audioContext.decoded, startSample, endSample, targetPoints);
    }
    
    function initChart() {
      if (!chartContainer || !audioContext.decoded) return;
      
      chart = echarts.init(chartContainer as HTMLElement);

      if (!chart) return;

      // Total number of stereo samples - used for fixed axis bounds
      const totalSamples = audioContext.decoded.length / 2;
      const initialData = getDataForZoom(0, 100);
      
      const option = {
        backgroundColor: chartColors.background,
        legend: createLegendConfig(),
        grid: createGridConfig(),
        tooltip: {
          ...createTooltipConfig(),
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
        xAxis: createXAxisConfig({ name: 'Sample', max: totalSamples }),
        yAxis: createYAxisConfig({ name: 'Amplitude', min: -1, max: 1 }),
        dataZoom: createDataZoomConfig(palette.left, palette.leftAlpha),
        series: [
          createLineSeriesConfig({
            name: 'Left Channel',
            data: initialData.left,
            color: palette.left,
            colorAlpha: palette.leftAlpha
          }),
          createLineSeriesConfig({
            name: 'Right Channel',
            data: initialData.right,
            color: palette.right,
            colorAlpha: palette.rightAlpha
          })
        ]
      };
      
      chart.setOption(option);
      
      // Listen to zoom events and update data resolution
      chart.on('dataZoom', () => {
        if (!chart) return;
        const { start, end } = getZoomFromOption(chart);
        
        // Only update if zoom changed significantly (avoid excessive updates)
        if (hasZoomChanged(currentZoom, start, end)) {
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
    
    const handleResize = createResizeHandler(() => chart);

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
