package handlers

import (
	"jiansuo/models"
	"jiansuo/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func Search(c *gin.Context) {
	var params services.SearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		Error(c, 400, err.Error())
		return
	}

	userIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	result, err := services.Search(params, userIP, userAgent)
	if err != nil {
		Error(c, 500, err.Error())
		return
	}

	Success(c, result)
}

func Suggest(c *gin.Context) {
	query := c.Query("q")
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 32)
	limit, _ := strconv.Atoi(c.Query("limit"))

	result := services.GetSuggestions(query, uint(categoryID), limit)
	Success(c, result)
}

func HotKeywords(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	result, err := services.GetHotKeywords(limit)
	if err != nil {
		Error(c, 500, err.Error())
		return
	}
	Success(c, result)
}

func IncrementView(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, 400, "invalid id")
		return
	}

	err = services.IncrementViewCount(uint(id))
	if err != nil {
		Error(c, 404, err.Error())
		return
	}

	Success(c, nil)
}

func KeywordDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, 400, "invalid id")
		return
	}

	kw, err := services.GetKeywordByID(uint(id))
	if err != nil {
		Error(c, 404, "keyword not found")
		return
	}

	Success(c, kw)
}

func CreateKeyword(c *gin.Context) {
	var kw models.Keyword
	if err := c.ShouldBindJSON(&kw); err != nil {
		Error(c, 400, err.Error())
		return
	}

	err := services.CreateKeyword(&kw)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, kw)
}

func UpdateKeyword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, 400, "invalid id")
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		Error(c, 400, err.Error())
		return
	}

	err = services.UpdateKeyword(uint(id), body)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, nil)
}

func DeleteKeyword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, 400, "invalid id")
		return
	}

	err = services.DeleteKeyword(uint(id))
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, nil)
}

func ListKeywords(c *gin.Context) {
	var params services.KeywordQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		Error(c, 400, err.Error())
		return
	}

	result, err := services.ListKeywords(params)
	if err != nil {
		Error(c, 500, err.Error())
		return
	}

	Success(c, result)
}

func BatchKeywords(c *gin.Context) {
	var body struct {
		Keywords []models.Keyword `json:"keywords"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		Error(c, 400, err.Error())
		return
	}

	count, err := services.BatchCreateKeywords(body.Keywords)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, gin.H{"imported_count": count})
}

func ImportCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		Error(c, 400, "file not found")
		return
	}
	defer file.Close()

	count, err := services.ImportCSV(file)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, gin.H{"imported_count": count})
}

func SetKeywordStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, 400, "invalid id")
		return
	}

	var body struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		Error(c, 400, err.Error())
		return
	}

	err = services.SetKeywordStatus(uint(id), body.Status)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, nil)
}

func CreateCategory(c *gin.Context) {
	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		Error(c, 400, err.Error())
		return
	}

	err := services.CreateCategory(&cat)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, cat)
}

func UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, 400, "invalid id")
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		Error(c, 400, err.Error())
		return
	}

	err = services.UpdateCategory(uint(id), body)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, nil)
}

func DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, 400, "invalid id")
		return
	}

	err = services.DeleteCategory(uint(id))
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, nil)
}

func ListCategories(c *gin.Context) {
	includeDisabled, _ := strconv.ParseBool(c.Query("all"))
	result, err := services.ListCategories(includeDisabled)
	if err != nil {
		Error(c, 500, err.Error())
		return
	}

	Success(c, result)
}

func SetCategoryStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, 400, "invalid id")
		return
	}

	var body struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		Error(c, 400, err.Error())
		return
	}

	err = services.SetCategoryStatus(uint(id), body.Status)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, nil)
}

func ListSearchLogs(c *gin.Context) {
	var params services.SearchLogQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		Error(c, 400, err.Error())
		return
	}

	result, err := services.ListSearchLogs(params)
	if err != nil {
		Error(c, 500, err.Error())
		return
	}

	Success(c, result)
}

func GetStatistics(c *gin.Context) {
	result, err := services.GetSearchStatistics()
	if err != nil {
		Error(c, 500, err.Error())
		return
	}

	Success(c, result)
}
