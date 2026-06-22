package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"whiteboard/internal/database"
	"whiteboard/internal/model"

	"github.com/gin-gonic/gin"
)

type WhiteboardHandler struct {
	WhiteboardRepo *database.WhiteboardRepo
	OperationRepo  *database.OperationRepo
	SnapshotRepo   *database.SnapshotRepo
}

func NewWhiteboardHandler(
	whiteboardRepo *database.WhiteboardRepo,
	operationRepo *database.OperationRepo,
	snapshotRepo *database.SnapshotRepo,
) *WhiteboardHandler {
	return &WhiteboardHandler{
		WhiteboardRepo: whiteboardRepo,
		OperationRepo:  operationRepo,
		SnapshotRepo:   snapshotRepo,
	}
}

type createWhiteboardRequest struct {
	Name       string `json:"name"`
	Background string `json:"background"`
}

func (h *WhiteboardHandler) CreateWhiteboard(c *gin.Context) {
	var req createWhiteboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	background := req.Background
	if background == "" {
		background = "#ffffff"
	}

	w := &model.Whiteboard{
		Name:       req.Name,
		Background: background,
	}

	if err := h.WhiteboardRepo.Create(w); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, w)
}

func (h *WhiteboardHandler) ListWhiteboards(c *gin.Context) {
	list, err := h.WhiteboardRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func parseID(c *gin.Context) (uint, bool) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return uint(id), true
}

type whiteboardDetailResponse struct {
	Whiteboard *model.Whiteboard `json:"whiteboard"`
	Operations []model.Operation `json:"operations"`
}

func (h *WhiteboardHandler) GetWhiteboard(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	w, err := h.WhiteboardRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if w == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "whiteboard not found"})
		return
	}

	ops, err := h.OperationRepo.ListByWhiteboard(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, whiteboardDetailResponse{
		Whiteboard: w,
		Operations: ops,
	})
}

type updateWhiteboardRequest struct {
	Name       string `json:"name"`
	Background string `json:"background"`
}

func (h *WhiteboardHandler) UpdateWhiteboard(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req updateWhiteboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	w, err := h.WhiteboardRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if w == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "whiteboard not found"})
		return
	}

	if req.Name != "" {
		w.Name = req.Name
	}
	if req.Background != "" {
		w.Background = req.Background
	}

	if err := h.WhiteboardRepo.Update(w); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, w)
}

func (h *WhiteboardHandler) DeleteWhiteboard(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	w, err := h.WhiteboardRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if w == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "whiteboard not found"})
		return
	}

	if err := h.OperationRepo.DeleteByWhiteboard(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.WhiteboardRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

type operationItem struct {
	UserID    uint   `json:"userId"`
	Type      string `json:"type"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

type saveOperationsRequest struct {
	Operations []operationItem `json:"operations"`
}

func (h *WhiteboardHandler) SaveOperations(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	w, err := h.WhiteboardRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if w == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "whiteboard not found"})
		return
	}

	var req saveOperationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ops := make([]model.Operation, 0, len(req.Operations))
	for _, item := range req.Operations {
		ops = append(ops, model.Operation{
			WhiteboardID: id,
			UserID:       item.UserID,
			Type:         item.Type,
			Data:         item.Data,
			Timestamp:    item.Timestamp,
		})
	}

	if err := h.OperationRepo.CreateBatch(ops); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "saved", "count": len(ops)})
}

type createSnapshotRequest struct {
	Name string `json:"name"`
}

func (h *WhiteboardHandler) CreateSnapshot(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	w, err := h.WhiteboardRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if w == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "whiteboard not found"})
		return
	}

	var req createSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := req.Name
	if name == "" {
		name = "Snapshot"
	}

	ops, err := h.OperationRepo.ListByWhiteboard(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	opsBytes, err := json.Marshal(ops)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s := &model.Snapshot{
		WhiteboardID: id,
		Name:         name,
		Operations:   string(opsBytes),
	}

	if err := h.SnapshotRepo.Create(s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, s)
}

func (h *WhiteboardHandler) ListSnapshots(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	w, err := h.WhiteboardRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if w == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "whiteboard not found"})
		return
	}

	list, err := h.SnapshotRepo.ListByWhiteboard(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *WhiteboardHandler) ClearWhiteboard(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	w, err := h.WhiteboardRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if w == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "whiteboard not found"})
		return
	}

	if err := h.OperationRepo.DeleteByWhiteboard(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cleared"})
}
