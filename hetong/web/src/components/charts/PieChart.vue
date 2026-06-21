<script setup lang="ts">
import { ref } from 'vue'
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
const { handlePieClick, getAnomalyItemStyle } = useDrillInteraction()

const getOption = (): EChartsOption => {
  const records = store.records
  const data = records.map((r, i) => ({
    value: r.metrics[props.metric] || 0,
    name: r.dimensions[props.dimension as any] || '未知',
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
      trigger: 'item',
      backgroundColor: 'rgba(15, 23, 42, 0.95)',
      borderColor: '#3B82F6',
      textStyle: { color: '#E2E8F0' },
      formatter: (params: any) => {
        const record = records[params.dataIndex]
        let html = `<div style="padding: 8px;">
          <div style="font-weight: bold; margin-bottom: 8px;">${params.name}</div>
          <div>${props.metric}: ${params.value.toLocaleString()}</div>
          <div>占比: ${params.percent}%</div>`

        if (record.comparison) {
          const yoyColor = record.comparison.yoyPercent >= 0 ? '#10B981' : '#EF4444'
          html += `<div style="color: ${yoyColor};">同比: ${record.comparison.yoyPercent.toFixed(2)}%</div>`
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
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      textStyle: { color: '#94A3B8' },
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['40%', '55%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 8,
          borderColor: '#0F172A',
          borderWidth: 2,
        },
        label: {
          show: false,
          position: 'center',
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 16,
            fontWeight: 'bold',
            color: '#E2E8F0',
          },
          itemStyle: {
            shadowBlur: 30,
            shadowColor: '#3B82F6',
          },
        },
        labelLine: {
          show: false,
        },
        data: data,
        color: [
          '#3B82F6', '#06B6D4', '#10B981', '#F59E0B',
          '#EF4444', '#8B5CF6', '#EC4899', '#14B8A6',
        ],
      },
    ],
  }
}

const { chartInstance } = useECharts(chartRef, getOption, [() => store.records])

chartInstance.value?.on('click', handlePieClick)
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
