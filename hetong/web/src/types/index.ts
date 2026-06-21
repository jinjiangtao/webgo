export type DimensionLevel = 'time' | 'region' | 'business'

export interface DrillLevel {
  dimension: DimensionLevel
  value: string
  label: string
}

export interface Filters {
  timeRange?: string[]
  regions?: string[]
  businessTypes?: string[]
  time?: string
  region?: string
  business?: string
  [key: string]: any
}

export interface AggregateRequest {
  dimensions: DimensionLevel[]
  metrics: string[]
  filters: Filters
  drillPath: DrillLevel[]
}

export interface Comparison {
  yoy: number
  mom: number
  yoyPercent: number
  momPercent: number
}

export interface Anomaly {
  isAnomaly: boolean
  severity: string
  reason: string
}

export interface AggRecord {
  dimensions: Record<DimensionLevel, string>
  metrics: Record<string, number>
  comparison?: Comparison
  anomaly?: Anomaly
  canDrillDown: boolean
  nextDimension: DimensionLevel
}

export interface RawRecord {
  id: number
  orderNo: string
  dimensions: Record<DimensionLevel, string>
  metrics: Record<string, number>
  createdAt: string
}

export interface AggregateResponse {
  code: number
  message: string
  data: {
    records: AggRecord[]
    summary: {
      totalRecords: number
      drillPath: DrillLevel[]
      availableDimensions: DimensionLevel[]
    }
  }
}

export interface DrillRequest {
  drillAction: 'down' | 'up'
  currentPath: DrillLevel[]
  drillDimension: string
  drillValue: string
  metrics: string[]
}

export interface DrillResponse {
  code: number
  message: string
  data: {
    drillPath: DrillLevel[]
    nextDimension: string
    records: AggRecord[]
  }
}

export interface TraceRequest {
  drillPath?: DrillLevel[]
  filters?: Filters
  page?: number
  pageSize?: number
  time?: string
  region?: string
  business?: string
  orderNo?: string
  [key: string]: any
}

export interface TraceResponse {
  code: number
  message: string
  data: {
    traceId: string
    aggregatePath: Array<Record<string, any>>
    records: RawRecord[]
    rawData: Array<Record<string, any>>
    total: number
    totalRaw: number
    page: number
    pageSize: number
  }
}

export interface DimensionResponse {
  code: number
  message: string
  data: {
    timeOptions: Array<{ value: string; label: string; children?: any[] }>
    regionOptions: Array<{ value: string; label: string; children?: any[] }>
    businessOptions: Array<{ value: string; label: string; children?: any[] }>
  }
}

export interface SnapshotItem {
  id: number
  name: string
  description?: string
  createdAt: string
  createdBy: string
  drillPath?: DrillLevel[]
  filters?: Filters
  records?: AggRecord[]
  recordCount?: number
}

export interface SnapshotRequest {
  name: string
  description?: string
  state?: Record<string, any>
  drillPath: DrillLevel[]
  filters: Filters
  records: AggRecord[]
  createdBy?: string
}

export interface ExportRequest {
  drillPath: DrillLevel[]
  filters: Filters
  metrics: string[]
  format?: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}
