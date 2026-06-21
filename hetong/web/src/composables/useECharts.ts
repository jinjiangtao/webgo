import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'

export function useECharts(
  chartRef: any,
  getOption: () => EChartsOption,
  deps?: any[]
) {
  const chartInstance = ref<echarts.ECharts | null>(null)

  const initChart = () => {
    if (!chartRef.value) return
    chartInstance.value = echarts.init(chartRef.value)
    updateChart()
  }

  const updateChart = () => {
    if (!chartInstance.value) return
    const option = getOption()
    chartInstance.value.setOption(option, true)
  }

  const resize = () => {
    chartInstance.value?.resize()
  }

  onMounted(() => {
    nextTick(() => {
      initChart()
      window.addEventListener('resize', resize)
    })
  })

  onUnmounted(() => {
    window.removeEventListener('resize', resize)
    chartInstance.value?.dispose()
  })

  if (deps) {
    watch(deps, () => {
      nextTick(() => updateChart())
    }, { deep: true })
  }

  return {
    chartInstance,
    initChart,
    updateChart,
    resize,
  }
}
