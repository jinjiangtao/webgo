package controller

import (
	"dashboard/internal/model"
	"dashboard/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DataController struct {
	aggregateService *service.AggregateService
	drillService     *service.DrillService
	traceService     *service.TraceService
	exportService    *service.ExportService
	snapshotService  *service.SnapshotService
}

func NewDataController() *DataController {
	return &DataController{
		aggregateService: service.NewAggregateService(),
		drillService:     service.NewDrillService(),
		traceService:     service.NewTraceService(),
		exportService:    service.NewExportService(),
		snapshotService:  service.NewSnapshotService(),
	}
}

func (c *DataController) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"status": "running",
		},
	})
}

func (c *DataController) Aggregate(ctx *gin.Context) {
	var req model.AggregateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.aggregateService.GetAggregatedData(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *DataController) Drill(ctx *gin.Context) {
	var req model.DrillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.drillService.ExecuteDrill(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "钻取失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *DataController) Trace(ctx *gin.Context) {
	var req model.TraceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	resp, err := c.traceService.GetTraceData(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "溯源查询失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *DataController) GetDimensions(ctx *gin.Context) {
	dimRepo := service.NewDimensionService()

	timeOpts, regionOpts, bizOpts, err := dimRepo.GetAllDimensions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取维度失败: " + err.Error(),
		})
		return
	}

	resp := &model.DimensionResponse{}
	resp.Code = 0
	resp.Message = "success"
	resp.Data.TimeOptions = timeOpts
	resp.Data.RegionOptions = regionOpts
	resp.Data.BusinessOptions = bizOpts

	ctx.JSON(http.StatusOK, resp)
}

func (c *DataController) SaveSnapshot(ctx *gin.Context) {
	var req model.SnapshotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	id, err := c.snapshotService.SaveSnapshot(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存快照失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"id": id,
		},
	})
}

func (c *DataController) GetSnapshots(ctx *gin.Context) {
	list, err := c.snapshotService.GetSnapshotList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取快照列表失败: " + err.Error(),
		})
		return
	}

	resp := &model.SnapshotListResponse{}
	resp.Code = 0
	resp.Message = "success"
	resp.Data = list

	ctx.JSON(http.StatusOK, resp)
}

func (c *DataController) GetSnapshot(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的ID",
		})
		return
	}
	id := uint(id64)

	detail, err := c.snapshotService.GetSnapshotDetail(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取快照详情失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, detail)
}

func (c *DataController) ExportExcel(ctx *gin.Context) {
	var req model.ExportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	filename, err := c.exportService.ExportExcel(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "导出失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"filename": filename,
		},
	})
}

func (c *DataController) ExportAggregateExcel(ctx *gin.Context) {
	var req model.ExportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	filename, err := c.exportService.ExportAggregateExcel(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "导出失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"filename": filename,
		},
	})
}
