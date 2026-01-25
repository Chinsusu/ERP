module github.com/erp-cosmetics/reporting-service

go 1.22

require (
	github.com/erp-cosmetics/shared v0.0.0
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/xuri/excelize/v2 v2.8.0
	go.uber.org/zap v1.26.0
	gorm.io/datatypes v1.2.0
	gorm.io/gorm v1.25.7
)

replace github.com/erp-cosmetics/shared => ../../shared
