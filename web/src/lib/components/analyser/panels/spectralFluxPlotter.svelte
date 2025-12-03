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
      downsampleMinMaxMono
    } from '$lib/utils/chartConfig';
    
    const audioContext = getAudioContext();
    const palette = chartPalettes.spectralFlux;
    
    let chartContainer = $state<HTMLElement | undefined>();
    let chart = $state<echarts.ECharts | null>(null);
    let currentZoom = $state({ start: 0, end: 100 });
    
    function getDataForZoom(start: number, end: number) {
      if (!audioContext.spectralFlux) return [];
      
      const totalFrames = audioContext.spectralFlux.length;
      const startIdx = Math.floor((start / 100) * totalFrames);
      const endIdx = Math.ceil((end / 100) * totalFrames);
      const visibleFrames = endIdx - startIdx;
      
      // Adjust target points based on zoom level
      const targetPoints = Math.min(visibleFrames, 700);
      
      return downsampleMinMaxMono(audioContext.spectralFlux, startIdx, endIdx, targetPoints);
    }
    
    function getMaxFluxValue(): number {
      if (!audioContext.spectralFlux) return 1;
      let max = 0;
      for (let i = 0; i < audioContext.spectralFlux.length; i++) {
        if (audioContext.spectralFlux[i] > max) {
          max = audioContext.spectralFlux[i];
        }
      }
      return max * 1.1; // Add 10% padding
    }
    
    function initChart() {
      if (!chartContainer || !audioContext.spectralFlux) return;
      
      chart = echarts.init(chartContainer as HTMLElement);

      if (!chart) return;

      const totalFrames = audioContext.spectralFlux.length;
      const maxValue = getMaxFluxValue();
      
      const option = {
        backgroundColor: chartColors.background,
        legend: createLegendConfig(),
        grid: createGridConfig({ leftGap: 60 }),
        tooltip: {
          ...createTooltipConfig(),
          formatter: (params: any) => {
            if (!params || params.length === 0) return '';
            const frame = params[0].data[0];
            let result = `Frame: ${frame}`;
            for (const param of params) {
              const value = param.data[1].toFixed(6);
              const color = param.color;
              result += `<br/><span style="color:${color}">●</span> Flux: ${value}`;
            }
            return result;
          }
        },
        xAxis: createXAxisConfig({ name: 'Frame', max: totalFrames }),
        yAxis: createYAxisConfig({
          name: 'Flux',
          min: 0,
          max: maxValue,
          nameGap: 50,
          labelFormatter: (value: number) => value.toExponential(1)
        }),
        dataZoom: createDataZoomConfig(palette.flux, palette.fluxAlpha),
        series: [
          createLineSeriesConfig({
            name: 'Spectral Flux',
            data: getDataForZoom(0, 100),
            color: palette.flux,
            colorAlpha: palette.fluxAlpha,
            lineWidth: 1.5,
            withGradient: true
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
              { data: newData }
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

<div class="spectral-flux-container" bind:this={chartContainer}></div>

<style>
    .spectral-flux-container {
      width: 100%;
      height: 100%;
      min-height: 0;
      background: var(--background);
      border-radius: var(--radius);
      border: 1px solid var(--border);
    }
</style>
