import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { DrillLevel, Filters, DimensionLevel, AggRecord } from '../types'
import { getAggregatedData, executeDrill, getDimensions } from '../api/data'

export const useDashboardStore = defineStore('dashboard', () => {
  const drillPath = ref<DrillLevel[]>([
    { dimension: 'time', value: '2025', label: '2025年' },
  ])
  const filters = ref<Filters>({})
  const metrics = ref<string[]>(['sales', 'orders', 'users'])
  const records = ref<AggRecord[]>([])
  const loading = ref(false)
  const dimensions = ref<{
    timeOptions: any[]
    regionOptions: any[]
    businessOptions: any[]
  }>({
    timeOptions: [],
    regionOptions: [],
    businessOptions: [],
  })

  const currentLevel = computed(() => drillPath.value.length)

  const canDrillDown = computed(() => {
    const usedDims = new Set(drillPath.value.map(d => d.dimension))
    return usedDims.size < 3
  })

  const pathSummary = computed(() => {
    if (drillPath.value.length === 0) return '全局视图'
    return drillPath.value.map(d => d.label).join(' → ')
  })

  const loadDimensions = async () => {
    try {
      const res = await getDimensions()
      dimensions.value = res.data
    } catch (e) {
      console.error('Failed to load dimensions:', e)
    }
  }

  const loadData = async () => {
    loading.value = true
    try {
      const usedDims = new Set(drillPath.value.map(d => d.dimension))
      const allDims: DimensionLevel[] = ['time', 'region', 'business']
      const availableDims = allDims.filter(d => !usedDims.has(d))

      const res = await getAggregatedData({
        dimensions: availableDims,
        metrics: metrics.value,
        filters: filters.value,
        drillPath: drillPath.value,
      })
      records.value = res.data.records
    } catch (e) {
      console.error('Failed to load data:', e)
    } finally {
      loading.value = false
    }
  }

  const drillDown = async (dimension: DimensionLevel, value: string, label: string) => {
    try {
      const res = await executeDrill({
        drillAction: 'down',
        currentPath: drillPath.value,
        drillDimension: dimension,
        drillValue: value,
        metrics: metrics.value,
      })
      drillPath.value = res.data.drillPath
      records.value = res.data.records
    } catch (e) {
      console.error('Failed to drill down:', e)
    }
  }

  const rollUp = async (targetIndex: number) => {
    if (targetIndex < 0 || targetIndex >= drillPath.value.length) return

    try {
      drillPath.value = drillPath.value.slice(0, targetIndex + 1)
      await loadData()
    } catch (e) {
      console.error('Failed to roll up:', e)
    }
  }

  const rollUpOne = async () => {
    if (drillPath.value.length <= 1) return
    await rollUp(drillPath.value.length - 2)
  }

  const updateFilters = (newFilters: Filters) => {
    filters.value = { ...filters.value, ...newFilters }
    loadData()
  }

  const updateMetrics = (newMetrics: string[]) => {
    metrics.value = newMetrics
    loadData()
  }

  return {
    drillPath,
    filters,
    metrics,
    records,
    loading,
    dimensions,
    currentLevel,
    canDrillDown,
    pathSummary,
    loadDimensions,
    loadData,
    drillDown,
    rollUp,
    rollUpOne,
    updateFilters,
    updateMetrics,
  }
})
