module github.com/erp-cosmetics/user-service

go 1.22

require (
	github.com/erp-cosmetics/shared v0.0.0
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.4.0
	github.com/redis/go-redis/v9 v9.3.0
	go.uber.org/zap v1.26.0
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
	gorm.io/gorm v1.25.5
)

replace github.com/erp-cosmetics/shared => ../../shared
