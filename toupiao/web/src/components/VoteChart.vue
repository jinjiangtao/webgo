<template>
  <div class="vote-chart">
    <div class="chart-header">
      <h3>
        <el-icon><PieChart /></el-icon>
        投票占比分布
      </h3>
      <el-radio-group v-model="chartType" size="small">
        <el-radio-button value="pie">饼图</el-radio-button>
        <el-radio-button value="bar">柱状图</el-radio-button>
      </el-radio-group>
    </div>
    <div ref="chartRef" class="chart-container"></div>
    <div class="chart-legend">
      <div
        v-for="(item, index) in chartData"
        :key="index"
        class="legend-item"
      >
        <span
          class="legend-color"
          :style="{ backgroundColor: colors[index % colors.length] }"
        ></span>
        <span class="legend-name">{{ item.name }}</span>
        <span class="legend-value">{{ item.value }} 票</span>
        <span class="legend-percent">{{ item.percentage.toFixed(1) }}%</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, onUnmounted, nextTick } from 'vue'
import * as echarts from 'echarts'

const props = defineProps({
  chartData: {
    type: Array,
    default: () => []
  }
})

const chartRef = ref(null)
const chartType = ref('pie')
let chartInstance = null

const colors = [
  '#667eea',
  '#764ba2',
  '#f093fb',
  '#f5576c',
  '#4facfe',
  '#43e97b',
  '#fa709a',
  '#fee140',
  '#30cfd0',
  '#a8edea'
]

const initChart = () => {
  if (!chartRef.value) return

  if (chartInstance) {
    chartInstance.dispose()
  }

  chartInstance = echarts.init(chartRef.value)
  updateChart()
}

const updateChart = () => {
  if (!chartInstance || props.chartData.length === 0) return

  const data = props.chartData.map((item, index) => ({
    ...item,
    itemStyle: {
      color: colors[index % colors.length]
    }
  }))

  let option = {}

  if (chartType.value === 'pie') {
    option = {
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} 票 ({d}%)'
      },
      series: [
        {
          name: '投票占比',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 8,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: true,
            formatter: '{b}\n{d}%'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 16,
              fontWeight: 'bold'
            },
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          },
          labelLine: {
            show: true
          },
          data: data
        }
      ]
    }
  } else {
    option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        },
        formatter: '{b}: {c} 票'
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: data.map(item => item.name),
        axisLabel: {
          interval: 0,
          rotate: 30
        }
      },
      yAxis: {
        type: 'value',
        name: '票数'
      },
      series: [
        {
          name: '票数',
          type: 'bar',
          barWidth: '60%',
          data: data.map(item => ({
            value: item.value,
            itemStyle: item.itemStyle
          })),
          label: {
            show: true,
            position: 'top',
            formatter: '{c}'
          }
        }
      ]
    }
  }

  chartInstance.setOption(option, true)
}

const handleResize = () => {
  chartInstance?.resize()
}

watch(() => props.chartData, () => {
  nextTick(() => {
    updateChart()
  })
}, { deep: true })

watch(chartType, () => {
  nextTick(() => {
    updateChart()
  })
})

onMounted(() => {
  nextTick(() => {
    initChart()
    window.addEventListener('resize', handleResize)
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
})
</script>

<style scoped>
.vote-chart {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.chart-header {
  padding: 20px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.chart-header h3 {
  margin: 0;
  font-size: 18px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.chart-container {
  height: 350px;
  padding: 20px;
}

.chart-legend {
  padding: 16px 24px;
  background: #f5f7fa;
  border-top: 1px solid #ebeef5;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.legend-name {
  color: #606266;
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.legend-value {
  font-weight: 600;
  color: #303133;
}

.legend-percent {
  color: #909399;
  font-size: 12px;
}
</style>
