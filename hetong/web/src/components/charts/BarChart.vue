<script setup lang="ts">
import { ref, computed } from 'vue'
import * as echarts from 'echarts'
import { useECharts } from '../../composables/useECharts'
import { useDrillInteraction } from '../../composables/useDrillInteraction'
import { useDashboardStore } from '../../stores/dashboard'
import type { EChartsOption } from 'echarts'

const props = defineProps<{
  title: string
  metric: string
  dimension: string
}>()

const chartRef = ref<HTMLDivElement | null>(null)
const store = useDashboardStore()
const { handleBarClick, getAnomalyItemStyle } = useDrillInteraction()

const getOption = (): EChartsOption => {
  const records = store.records
  const xAxisData = records.map(r => r.dimensions[props.dimension as any] || '未知')
  const seriesData = records.map((r, i) => ({
    value: r.metrics[props.metric] || 0,
    ...getAnomalyItemStyle(i),
  }))

  return {
    backgroundColor: 'transparent',
    title: {
      text: props.title,
      textStyle: {
        color: '#E2E8F0',
        fontSize: 14,
        fontFamily: 'JetBrains Mono, monospace',
      },
      left: 'center',
      top: 10,
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(15, 23, 42, 0.95)',
      borderColor: '#3B82F6',
      textStyle: { color: '#E2E8F0' },
      formatter: (params: any) => {
        const data = params[0]
        const record = records[data.dataIndex]
        let html = `<div style="padding: 8px;">
          <div style="font-weight: bold; margin-bottom: 8px;">${data.name}</div>
          <div>${props.metric}: ${data.value.toLocaleString()}</div>`

        if (record.comparison) {
          const yoyColor = record.comparison.yoyPercent >= 0 ? '#10B981' : '#EF4444'
          const momColor = record.comparison.momPercent >= 0 ? '#10B981' : '#EF4444'
          html += `<div style="color: ${yoyColor};">同比: ${record.comparison.yoyPercent.toFixed(2)}%</div>`
          html += `<div style="color: ${momColor};">环比: ${record.comparison.momPercent.toFixed(2)}%</div>`
        }

        if (record.anomaly?.isAnomaly) {
          html += `<div style="color: #EF4444; font-weight: bold; margin-top: 4px;">⚠ ${record.anomaly.reason}</div>`
        }

        if (record.canDrillDown) {
          html += `<div style="color: #60A5FA; margin-top: 8px;">点击下钻 →</div>`
        }
        html += '</div>'
        return html
      },
    },
    grid: {
      left: '10%',
      right: '10%',
      top: '15%',
      bottom: '15%',
    },
    xAxis: {
      type: 'category',
      data: xAxisData,
      axisLine: { lineStyle: { color: '#475569' } },
      axisLabel: { color: '#94A3B8', rotate: 30 },
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisLabel: { color: '#94A3B8' },
      splitLine: { lineStyle: { color: '#1E293B' } },
    },
    series: [
      {
        type: 'bar',
        data: seriesData,
        barWidth: '50%',
        itemStyle: {
          borderRadius: [4, 4, 0, 0],
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#3B82F6' },
            { offset: 1, color: '#06B6D4' },
          ]),
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 20,
            shadowColor: '#3B82F6',
          },
        },
      },
    ],
  }
}

const { chartInstance } = useECharts(chartRef, getOption, [() => store.records])

chartInstance.value?.on('click', handleBarClick)
</script>

<template>
  <div class="chart-wrapper">
    <div ref="chartRef" class="chart-container"></div>
  </div>
</template>

<style scoped>
.chart-wrapper {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.8), rgba(15, 23, 42, 0.9));
  border: 1px solid rgba(59, 130, 246, 0.3);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.chart-wrapper:hover {
  border-color: rgba(59, 130, 246, 0.6);
  box-shadow: 0 0 30px rgba(59, 130, 246, 0.2);
}

.chart-container {
  width: 100%;
  height: 100%;
  min-height: 300px;
}
</style>
