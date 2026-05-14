package router

import (
	"github.com/friedrichad/golang_web_api_demo/internal/controller"
	"github.com/friedrichad/golang_web_api_demo/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/uploads", "./internal/upload")
	configCors(router)
	initAuthRouter(router)
	initUserRouter(router)
	initBinRouter(router)
	initWarehouseRouter(router)
	initComponentRouter(router)
	initCustomerRouter(router)
	initRoleRouter(router)
	initPositionRouter(router)
	initRequestRouter(router)
	initRequestDetailRouter(router)
	initRequestPermissionRouter(router)
	initInventoryAdjustmentRouter(router)
	initInventoryAuditRouter(router)
	initInventoryAuditDetailRouter(router)
	initInventoryLedgerRouter(router)
	initComponentCategoryRouter(router)
	initUploadRouter(router)
	return router
}

func configCors(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:  viper.GetStringSlice("security.cors"),
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}))
}

func initUserRouter(router *gin.Engine) {
	userController := controller.NewUserController()
	userGroup := router.Group("/users")
	userGroup.Use(middleware.BearerAuthenticator())
	{
		userGroup.GET("", middleware.Authorizator("user:view"), userController.GetAllUsers())
		userGroup.GET("/:id", middleware.Authorizator("user:view"), userController.GetUserById())
		userGroup.POST("", middleware.Authorizator("user:create"), userController.CreateUser())
		userGroup.PUT("", middleware.Authorizator("user:edit"), userController.UpdateUser())
		userGroup.DELETE("", middleware.Authorizator("user:delete"), userController.DeleteUser())
		userGroup.GET("/:id/authorities", middleware.Authorizator("user:view"), userController.GetUserAuthorities())
	}
}

func initBinRouter(router *gin.Engine) {
	binController := controller.NewBinController()
	binGroup := router.Group("/bins")
	binGroup.Use(middleware.BearerAuthenticator())
	{
		binGroup.GET("", middleware.Authorizator("bin:view"), binController.GetAllBins())
		binGroup.GET("/:id", middleware.Authorizator("bin:view"), binController.GetBinById())
		binGroup.POST("", middleware.Authorizator("bin:create"), binController.CreateBin())
		binGroup.PUT("", middleware.Authorizator("bin:edit"), binController.UpdateBin())
		binGroup.DELETE("", middleware.Authorizator("bin:delete"), binController.DeleteBin())
	}
}

func initWarehouseRouter(router *gin.Engine) {
	warehouseController := controller.NewWarehouseController()
	warehouseGroup := router.Group("/warehouses")
	warehouseGroup.Use(middleware.BearerAuthenticator())
	{
		warehouseGroup.GET("", middleware.Authorizator("warehouse:view"), warehouseController.GetAllWarehouses())
		warehouseGroup.GET("/:id", middleware.Authorizator("warehouse:view"), warehouseController.GetWarehouseById())
		warehouseGroup.POST("", middleware.Authorizator("warehouse:create"), warehouseController.CreateWarehouse())
		warehouseGroup.PUT("", middleware.Authorizator("warehouse:edit"), warehouseController.UpdateWarehouse())
		warehouseGroup.DELETE("", middleware.Authorizator("warehouse:delete"), warehouseController.DeleteWarehouse())
	}
}

func initComponentRouter(router *gin.Engine) {
	componentController := controller.NewComponentController()
	componentGroup := router.Group("/components")
	componentGroup.Use(middleware.BearerAuthenticator())
	{
		componentGroup.GET("", middleware.Authorizator("component:view"), componentController.GetAllComponents())
		componentGroup.GET("/:id", middleware.Authorizator("component:view"), componentController.GetComponentById())
		componentGroup.POST("", middleware.Authorizator("component:create"), componentController.CreateComponent())
		componentGroup.PUT("", middleware.Authorizator("component:edit"), componentController.UpdateComponent())
		componentGroup.DELETE("", middleware.Authorizator("component:delete"), componentController.DeleteComponent())
	}
}

func initCustomerRouter(router *gin.Engine) {
	customerController := controller.NewCustomerController()
	customerGroup := router.Group("/customers")
	customerGroup.Use(middleware.BearerAuthenticator())
	{
		customerGroup.GET("", middleware.Authorizator("customer:view"), customerController.GetAllCustomers())
		customerGroup.GET("/:id", middleware.Authorizator("customer:view"), customerController.GetCustomerById())
		customerGroup.POST("", middleware.Authorizator("customer:create"), customerController.CreateCustomer())
		customerGroup.PUT("", middleware.Authorizator("customer:edit"), customerController.UpdateCustomer())
		customerGroup.DELETE("", middleware.Authorizator("customer:delete"), customerController.DeleteCustomer())
	}
}

func initRoleRouter(router *gin.Engine) {
	roleController := controller.NewRoleController()
	roleGroup := router.Group("/roles")
	roleGroup.Use(middleware.BearerAuthenticator())
	{
		roleGroup.GET("", middleware.Authorizator("role:view"), roleController.GetAllRoles())
		roleGroup.GET("/:id", middleware.Authorizator("role:view"), roleController.GetRoleById())
		roleGroup.POST("", middleware.Authorizator("role:create"), roleController.CreateRole())
		roleGroup.PUT("", middleware.Authorizator("role:edit"), roleController.UpdateRole())
		roleGroup.DELETE("", middleware.Authorizator("role:delete"), roleController.DeleteRole())
	}
}

func initPositionRouter(router *gin.Engine) {
	positionController := controller.NewPositionController()
	positionGroup := router.Group("/positions")
	positionGroup.Use(middleware.BearerAuthenticator())
	{
		positionGroup.GET("", middleware.Authorizator("position:view"), positionController.GetAllPositions())
		positionGroup.GET("/:id", middleware.Authorizator("position:view"), positionController.GetPositionById())
		positionGroup.POST("", middleware.Authorizator("position:create"), positionController.CreatePosition())
		positionGroup.PUT("", middleware.Authorizator("position:edit"), positionController.UpdatePosition())
		positionGroup.DELETE("", middleware.Authorizator("position:delete"), positionController.DeletePosition())
	}
}

func initRequestRouter(router *gin.Engine) {
	requestController := controller.NewRequestController()
	requestGroup := router.Group("/requests")
	requestGroup.Use(middleware.BearerAuthenticator())
	{
		requestGroup.GET("", middleware.Authorizator("request:view"), requestController.GetAllRequests())
		requestGroup.GET("/:id", middleware.Authorizator("request:view"), requestController.GetRequestById())
		requestGroup.POST("", middleware.Authorizator("request:create"), requestController.CreateRequest())
		requestGroup.PUT("", middleware.Authorizator("request:edit"), requestController.UpdateRequest())
		requestGroup.DELETE("", middleware.Authorizator("request:delete"), requestController.DeleteRequest())
		requestGroup.POST("/approval", middleware.Authorizator("request:approve"), requestController.ApprovalRequest())
		requestGroup.POST("/confirm", middleware.Authorizator("request:confirm"), requestController.ConfirmRequest())
	}
}

func initRequestDetailRouter(router *gin.Engine) {
	requestDetailController := controller.NewRequestDetailController()
	requestDetailGroup := router.Group("/request-details")
	requestDetailGroup.Use(middleware.BearerAuthenticator())
	{
		requestDetailGroup.GET("", middleware.Authorizator("request:view"), requestDetailController.GetAllRequestDetails())
		requestDetailGroup.GET("/:id", middleware.Authorizator("request:view"), requestDetailController.GetRequestDetailById())
		requestDetailGroup.POST("", middleware.Authorizator("request:create"), requestDetailController.CreateRequestDetail())
		requestDetailGroup.PUT("", middleware.Authorizator("request:edit"), requestDetailController.UpdateRequestDetail())
		requestDetailGroup.DELETE("", middleware.Authorizator("request:delete"), requestDetailController.DeleteRequestDetail())
	}
}

func initRequestPermissionRouter(router *gin.Engine) {
	permissionController := controller.NewRequestPermissionController()
	permissionGroup := router.Group("/request-permissions")
	permissionGroup.Use(middleware.BearerAuthenticator())
	{
		permissionGroup.GET("", middleware.Authorizator("request:view"), permissionController.GetAllPermissions())
		permissionGroup.POST("", middleware.Authorizator("request:create"), permissionController.CreatePermission())
		permissionGroup.PUT("", middleware.Authorizator("request:edit"), permissionController.UpdatePermission())
		permissionGroup.DELETE("", middleware.Authorizator("request:delete"), permissionController.DeletePermission())
		permissionGroup.POST("/approval", middleware.Authorizator("request:approve"), permissionController.ApprovalPermission())
	}
}

func initInventoryAdjustmentRouter(router *gin.Engine) {
	adjustmentController := controller.NewInventoryAdjustmentController()
	adjustmentGroup := router.Group("/adjustments")
	adjustmentGroup.Use(middleware.BearerAuthenticator())
	{
		adjustmentGroup.GET("", middleware.Authorizator("adjustment:view"), adjustmentController.GetAllAdjustments())
		adjustmentGroup.GET("/:id", middleware.Authorizator("adjustment:view"), adjustmentController.GetAdjustmentById())
		adjustmentGroup.POST("", middleware.Authorizator("adjustment:create"), adjustmentController.CreateAdjustment())
		adjustmentGroup.PUT("", middleware.Authorizator("adjustment:edit"), adjustmentController.UpdateAdjustment())
		adjustmentGroup.DELETE("", middleware.Authorizator("adjustment:delete"), adjustmentController.DeleteAdjustment())
		adjustmentGroup.POST("/approval", middleware.Authorizator("adjustment:approve"), adjustmentController.ApproveAdjustment())
	}
}

func initInventoryAuditRouter(router *gin.Engine) {
	auditController := controller.NewInventoryAuditController()
	auditGroup := router.Group("/audits")
	auditGroup.Use(middleware.BearerAuthenticator())
	{
		auditGroup.GET("", middleware.Authorizator("audit:view"), auditController.GetAllAudits())
		auditGroup.GET("/:id", middleware.Authorizator("audit:view"), auditController.GetAuditById())
		auditGroup.POST("", middleware.Authorizator("audit:create"), auditController.CreateAudit())
		auditGroup.PUT("", middleware.Authorizator("audit:edit"), auditController.UpdateAudit())
		auditGroup.DELETE("", middleware.Authorizator("audit:delete"), auditController.DeleteAudit())
		auditGroup.POST("/approval", middleware.Authorizator("audit:approve"), auditController.ApproveAudit())
		auditGroup.POST("/confirm", middleware.Authorizator("audit:confirm"), auditController.ConfirmAudit())
	}
}
func initInventoryAuditDetailRouter(router *gin.Engine) {
	auditDetailController := controller.NewInventoryAuditDetailController()
	auditDetailGroup := router.Group("/audit-details")
	auditDetailGroup.Use(middleware.BearerAuthenticator())
	{
		auditDetailGroup.GET("", middleware.Authorizator("audit:view"), auditDetailController.GetAllInventoryAuditDetails())
		// auditDetailGroup.GET("/:id", auditDetailController.GetInventoryAuditDetailById())
		auditDetailGroup.POST("", middleware.Authorizator("audit:create"), auditDetailController.CreateInventoryAuditDetail())
		auditDetailGroup.PUT("", middleware.Authorizator("audit:edit"), auditDetailController.UpdateInventoryAuditDetail())
		auditDetailGroup.DELETE("", middleware.Authorizator("audit:delete"), auditDetailController.DeleteInventoryAuditDetail())
	}
}

func initInventoryLedgerRouter(router *gin.Engine) {
	ledgerController := controller.NewInventoryLedgerController()
	ledgerGroup := router.Group("/ledgers")
	ledgerGroup.Use(middleware.BearerAuthenticator())
	{
		ledgerGroup.GET("", middleware.Authorizator("ledgers:view"), ledgerController.GetAllLedgers())
		ledgerGroup.GET("/export", middleware.Authorizator("ledgers:view"), ledgerController.ExportLedgersExcel())
		ledgerGroup.GET("/:id", middleware.Authorizator("ledgers:view"), ledgerController.GetLedgerById())
	}
}

func initComponentCategoryRouter(router *gin.Engine) {
	categoryController := controller.NewComponentCategoryController()
	categoryGroup := router.Group("/categories")
	categoryGroup.Use(middleware.BearerAuthenticator())
	{
		categoryGroup.GET("", middleware.Authorizator("categories:view"), categoryController.GetAllCategories())
		categoryGroup.GET("/:id", middleware.Authorizator("categories:view"), categoryController.GetCategoryById())
		categoryGroup.POST("", middleware.Authorizator("categories:create"), categoryController.CreateCategory())
		categoryGroup.PUT("", middleware.Authorizator("categories:edit"), categoryController.UpdateCategory())
		categoryGroup.DELETE("", middleware.Authorizator("categories:delete"), categoryController.DeleteCategory())
	}
}

func initAuthRouter(router *gin.Engine) {
	authController := controller.NewAuthController()
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/refresh", authController.GetToken())
		authGroup.POST("/login", authController.GetToken())
		authGroup.POST("/register", authController.Register())
		authGroup.POST("/logout", middleware.BearerAuthenticator(), authController.Logout())
	}
}

func initUploadRouter(router *gin.Engine) {
	uploadController := controller.NewUploadController()
	uploadGroup := router.Group("/uploads")
	uploadGroup.Use(middleware.BearerAuthenticator())
	{
		uploadGroup.POST("/base64", middleware.Authorizator("upload:create"), uploadController.UploadBase64())
		uploadGroup.POST("/multipart", middleware.Authorizator("upload:create"), uploadController.UploadMultipart())
		uploadGroup.POST("/multiple", middleware.Authorizator("upload:create"), uploadController.UploadMultiple())
	}
}
