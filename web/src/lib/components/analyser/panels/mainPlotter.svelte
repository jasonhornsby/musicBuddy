<script lang="ts">
	import { getAudioContext } from '$lib/context/audio.svelte';
    import * as echarts from 'echarts';
    import { onMount } from "svelte";

    let chartContainer = $state<HTMLElement | undefined>();
    let chart = $state<echarts.ECharts | null>(null);

    const audioContext = getAudioContext();
    const bufferView = audioContext.getAudioBufferView();

    function getDownsampledChannelData(targetPoints: number, slice: { start: number, end: number } | null = null) {
        const downSampledChannelData: Float32Array[] = [];
        for (let i = 0; i < bufferView.getChannelViews().length; i++) {
            downSampledChannelData.push(...bufferView.getDownsampledMinMax(i, targetPoints, slice));
        }
        return downSampledChannelData;
    }

    function initChart() {
        if (!chartContainer) return;

        chart = echarts.init(chartContainer);

        const targetPoints = 700;
        const downSampledChannelData = getDownsampledChannelData(targetPoints);
        const option: echarts.EChartsOption = {
            tooltip: {
                trigger: 'axis',
                axisPointer: {
                    type: 'cross'
                },
                formatter: (params: any) => {
                    if (!Array.isArray(params) || params.length === 0) return '';
                    const frameNumber = params[0].value[0];
                    const seconds = (frameNumber / bufferView.getAudioInfo().sampleRate).toFixed(4);
                    let result = `<strong>${seconds}s (sample ${frameNumber})</strong>`;
                    params.forEach((param: any) => {
                        const min = param.value[1]?.toFixed(4) ?? 'N/A';
                        const max = param.value[2]?.toFixed(4) ?? 'N/A';
                        result += `<br/>${param.marker} ${param.seriesName}: ${min} ~ ${max}`;
                    });
                    return result;
                }
            },
            legend: {
                show: true,
                top: 0,
                left: 'center',
                // Only show the actual channel series, not the preview series
                data: bufferView.getChannelViews().map((_, index) => `Channel ${index + 1}`)
            },
            grid: {
                left: 40,
                right: 10,
                top: 30,
                bottom: 60,
                containLabel: false
            },
            dataZoom: [
                {
                    type: 'slider',
                    xAxisIndex: 0,
                    start: 0,
                    end: 100,
                    bottom: 10,
                    showDataShadow: true,
                },
                {
                    type: 'inside',
                    xAxisIndex: 0,
                    start: 0,
                    end: 100,
                }
            ],
            dataset: {
                source: downSampledChannelData as any,
            },
            xAxis: { type: 'value', min: 0, max: bufferView.numSamples, axisLabel: {
                formatter: (frameNumber: number) => {
                    return `${(frameNumber / bufferView.getAudioInfo().sampleRate).toFixed(2)}s`;
                }
            } },
            yAxis: { type: 'value', min: -1, max: 1 },
            series: [
                // Hidden line series for dataZoom preview (one per channel)
                ...bufferView.getChannelViews().map(
                    (_, index) => ({
                        type: 'line' as const,
                        name: `Channel ${index + 1} Preview`,
                        seriesLayoutBy: 'row' as const,
                        animation: false,
                        showSymbol: false,
                        silent: true,
                        large: true,
                        lineStyle: { opacity: 0 },
                        areaStyle: { opacity: 0 },
                        encode: {
                            x: index * 3,
                            y: index * 3 + 2, // Use max value for preview shape
                        },
                        // Hide from legend and tooltip
                        legendHoverLink: false,
                        tooltip: { show: false },
                    })
                ),
                // Visible custom series for actual waveform rendering
                ...bufferView.getChannelViews().map(
                    (_, index) => ({
                        type: 'custom' as const,
                        name: `Channel ${index + 1}`,
                        seriesLayoutBy: 'row' as const,
                        animation: false,
                        large: true,
                        encode: {
                            x: index * 3,
                            y: [index * 3 + 1, index * 3 + 2],
                        },
                        renderItem: (params: any, api: any) => {
                            const xValue = api.value(0);
                            const start = api.coord([xValue, api.value(1)]);
                            const end = api.coord([xValue, api.value(2)]);
                            return {
                                type: 'line' as const,
                                shape: {
                                    x1: start[0], y1: start[1],
                                    x2: end[0], y2: end[1]
                                },
                                style: api.style({
                                    stroke: api.visual('color'),
                                    lineWidth: 0.5
                                })
                            }
                        }
                    })
                )
            ]
        }
        chart.setOption(option);

        chart.on('dataZoom', () => {
            if (!chart) return;
            const option = chart.getOption() as echarts.EChartsOption;
            const dataZoom = option.dataZoom;
            if (!dataZoom || !Array.isArray(dataZoom)) return;

            const start = Math.floor(dataZoom[0].startValue as number || 0);
            const end = Math.ceil(dataZoom[0].endValue as number || bufferView.numSamples);

            const downSampledChannelData = getDownsampledChannelData(targetPoints, { start, end });
            chart.setOption({
                dataset: {
                    source: downSampledChannelData as any
                }
            });
        })

        window.addEventListener('resize', () => chart?.resize());
    }

    onMount(() => {
        initChart();
        return () => {
            chart?.dispose();
            window.removeEventListener('resize', () => chart?.resize());
        };
    });
</script>

<div class="w-full h-full min-h-0" bind:this={chartContainer}></div>
