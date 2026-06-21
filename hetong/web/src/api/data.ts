import { get, post, postBlob } from './request'
import type {
  AggregateRequest,
  AggregateResponse,
  DrillRequest,
  DrillResponse,
  TraceRequest,
  TraceResponse,
  DimensionResponse,
  SnapshotRequest,
  SnapshotItem,
  ExportRequest,
  ApiResponse,
} from '../types'

export const getHealth = (): Promise<ApiResponse> => {
  return get('/health')
}

export const getDimensions = (): Promise<DimensionResponse> => {
  return get('/dimensions')
}

export const getAggregatedData = (params: AggregateRequest): Promise<AggregateResponse> => {
  return post('/data/aggregate', params)
}

export const executeDrill = (params: DrillRequest): Promise<DrillResponse> => {
  return post('/data/drill', params)
}

export const getTraceData = (params: TraceRequest): Promise<TraceResponse> => {
  return post('/data/trace', params)
}

export const saveSnapshot = (params: SnapshotRequest): Promise<ApiResponse> => {
  return post('/snapshots', params)
}

export const getSnapshotList = (): Promise<ApiResponse<SnapshotItem[]>> => {
  return get('/snapshots')
}

export const getSnapshotDetail = (id: number): Promise<ApiResponse<SnapshotItem>> => {
  return get(`/snapshots/${id}`)
}

export const exportExcel = (params: ExportRequest): Promise<Blob> => {
  return postBlob('/export/excel', params)
}

export const exportAggregateExcel = (params: ExportRequest): Promise<Blob> => {
  return postBlob('/export/aggregate-excel', params)
}
