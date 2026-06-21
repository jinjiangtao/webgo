import { useDashboardStore } from '../stores/dashboard'
import type { DimensionLevel } from '../types'

export function useDrillInteraction() {
  const store = useDashboardStore()

  const handleBarClick = (params: any) => {
    if (!store.canDrillDown) return

    const nextDim = getNextDimension()
    if (!nextDim) return

    const value = params.name
    const label = value
    store.drillDown(nextDim, value, label)
  }

  const handlePieClick = (params: any) => {
    if (!store.canDrillDown) return

    const nextDim = getNextDimension()
    if (!nextDim) return

    const value = params.name
    const label = value
    store.drillDown(nextDim, value, label)
  }

  const handleMapClick = (params: any) => {
    if (!store.canDrillDown) return

    const nextDim = getNextDimension()
    if (!nextDim) return

    const value = params.name
    const label = value
    store.drillDown(nextDim, value, label)
  }

  const getNextDimension = (): DimensionLevel | null => {
    const usedDims = new Set(store.drillPath.map(d => d.dimension))
    const allDims: DimensionLevel[] = ['time', 'region', 'business']

    for (const dim of allDims) {
      if (!usedDims.has(dim)) {
        return dim
      }
    }
    return null
  }

  const getAnomalyItems = () => {
    return store.records
      .map((r, i) => ({ index: i, ...r }))
      .filter(r => r.anomaly?.isAnomaly)
  }

  const getAnomalyItemStyle = (index: number) => {
    const record = store.records[index]
    if (!record?.anomaly?.isAnomaly) return {}

    const severity = record.anomaly.severity
    const colors: Record<string, string> = {
      low: '#F59E0B',
      medium: '#F97316',
      high: '#EF4444',
    }

    return {
      itemStyle: {
        color: colors[severity] || '#EF4444',
        shadowBlur: 20,
        shadowColor: colors[severity] || '#EF4444',
      },
    }
  }

  return {
    handleBarClick,
    handlePieClick,
    handleMapClick,
    getNextDimension,
    getAnomalyItems,
    getAnomalyItemStyle,
  }
}
