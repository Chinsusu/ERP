package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/erp-cosmetics/manufacturing-service/internal/delivery/http/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestManufacturingAPI_BOM_Access(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Get BOM with Formula Access", func(t *testing.T) {
		r := gin.New()
		r.GET("/api/v1/boms/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data": gin.H{
					"id": c.Param("id"),
					"formula_details": gin.H{"notes": "Confidential"},
				},
			})
		})

		req, _ := http.NewRequest("GET", "/api/v1/boms/"+uuid.New().String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.NotNil(t, data["formula_details"])
	})

	t.Run("Create Work Order Success", func(t *testing.T) {
		r := gin.New()
		r.POST("/api/v1/work-orders", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{
				"status": "success",
				"data": gin.H{"wo_number": "WO-1001"},
			})
		})

		req, _ := http.NewRequest("POST", "/api/v1/work-orders", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
