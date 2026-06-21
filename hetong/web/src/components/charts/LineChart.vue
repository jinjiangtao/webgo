<script setup lang="ts">
import { ref } from 'vue'
import * as echarts from 'echarts'
import { useECharts } from '../../composables/useECharts'
import { useDashboardStore } from '../../stores/dashboard'
import type { EChartsOption } from 'echarts'

const props = defineProps<{
  title: string
  metric: string
  dimension: string
}>()

const chartRef = ref<HTMLDivElement | null>(null)
const store = useDashboardStore()

const getOption = (): EChartsOption => {
  const records = store.records
  const xAxisData = records.map(r => r.dimensions[props.dimension as any] || '未知')
  const seriesData = records.map(r => r.metrics[props.metric] || 0)
  const yoyData = records.map(r => r.comparison?.yoyPercent || 0)

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
    },
    legend: {
      data: [props.metric, '同比'],
      textStyle: { color: '#94A3B8' },
      top: 35,
    },
    grid: {
      left: '10%',
      right: '10%',
      top: '20%',
      bottom: '15%',
    },
    xAxis: {
      type: 'category',
      data: xAxisData,
      axisLine: { lineStyle: { color: '#475569' } },
      axisLabel: { color: '#94A3B8', rotate: 30 },
    },
    yAxis: [
      {
        type: 'value',
        axisLine: { show: false },
        axisLabel: { color: '#94A3B8' },
        splitLine: { lineStyle: { color: '#1E293B' } },
      },
      {
        type: 'value',
        axisLine: { show: false },
        axisLabel: { color: '#10B981', formatter: '{value}%' },
        splitLine: { show: false },
      },
    ],
    series: [
      {
        name: props.metric,
        type: 'line',
        smooth: true,
        data: seriesData,
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(59, 130, 246, 0.4)' },
            { offset: 1, color: 'rgba(59, 130, 246, 0.05)' },
          ]),
        },
        lineStyle: {
          width: 3,
          color: '#3B82F6',
        },
        itemStyle: {
          color: '#3B82F6',
        },
      },
      {
        name: '同比',
        type: 'line',
        smooth: true,
        yAxisIndex: 1,
        data: yoyData,
        lineStyle: {
          width: 2,
          color: '#10B981',
          type: 'dashed',
        },
        itemStyle: {
          color: '#10B981',
        },
      },
    ],
  }
}

useECharts(chartRef, getOption, [() => store.records])
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
