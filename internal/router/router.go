package router

import (
	"github.com/friedrichad/golang_web_api_demo/internal/controller"
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
	initRequestRouter(router)
	initRequestDetailRouter(router)
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
	{
		userGroup.GET("", userController.GetAllUsers())
		userGroup.GET("/:id", userController.GetUserById())
		userGroup.POST("", userController.CreateUser())
		userGroup.PUT("", userController.UpdateUser())
		userGroup.DELETE("", userController.DeleteUser())
		userGroup.GET("/:id/authorities", userController.GetUserAuthorities())
	}
}

func initBinRouter(router *gin.Engine) {
	binController := controller.NewBinController()
	binGroup := router.Group("/bins")
	{
		binGroup.GET("", binController.GetAllBins())
		binGroup.GET("/:id", binController.GetBinById())
		binGroup.POST("", binController.CreateBin())
		binGroup.PUT("", binController.UpdateBin())
		binGroup.DELETE("", binController.DeleteBin())
	}
}

func initWarehouseRouter(router *gin.Engine) {
	warehouseController := controller.NewWarehouseController()
	warehouseGroup := router.Group("/warehouses")
	{
		warehouseGroup.GET("", warehouseController.GetAllWarehouses())
		warehouseGroup.GET("/:id", warehouseController.GetWarehouseById())
		warehouseGroup.POST("", warehouseController.CreateWarehouse())
		warehouseGroup.PUT("", warehouseController.UpdateWarehouse())
		warehouseGroup.DELETE("", warehouseController.DeleteWarehouse())
	}
}

func initComponentRouter(router *gin.Engine) {
	componentController := controller.NewComponentController()
	componentGroup := router.Group("/components")
	{
		componentGroup.GET("", componentController.GetAllComponents())
		componentGroup.GET("/:id", componentController.GetComponentById())
		componentGroup.POST("", componentController.CreateComponent())
		componentGroup.PUT("", componentController.UpdateComponent())
		componentGroup.DELETE("", componentController.DeleteComponent())
	}
}

func initCustomerRouter(router *gin.Engine) {
	customerController := controller.NewCustomerController()
	customerGroup := router.Group("/customers")
	{
		customerGroup.GET("", customerController.GetAllCustomers())
		customerGroup.GET("/:id", customerController.GetCustomerById())
		customerGroup.POST("", customerController.CreateCustomer())
		customerGroup.PUT("", customerController.UpdateCustomer())
		customerGroup.DELETE("", customerController.DeleteCustomer())
	}
}

func initRoleRouter(router *gin.Engine) {
	roleController := controller.NewRoleController()
	roleGroup := router.Group("/roles")
	{
		roleGroup.GET("", roleController.GetAllRoles())
		roleGroup.GET("/:id", roleController.GetRoleById())
		roleGroup.POST("", roleController.CreateRole())
		roleGroup.PUT("", roleController.UpdateRole())
		roleGroup.DELETE("", roleController.DeleteRole())
	}
}

func initRequestRouter(router *gin.Engine) {
	requestController := controller.NewRequestController()
	requestGroup := router.Group("/requests")
	{
		requestGroup.GET("", requestController.GetAllRequests())
		requestGroup.GET("/:id", requestController.GetRequestById())
		requestGroup.POST("", requestController.CreateRequest())
		requestGroup.PUT("", requestController.UpdateRequest())
		requestGroup.DELETE("", requestController.DeleteRequest())
		requestGroup.POST("/approval", requestController.ApprovalRequest())
		requestGroup.POST("/confirm", requestController.ConfirmRequest())
	}
}

func initRequestDetailRouter(router *gin.Engine) {
	requestDetailController := controller.NewRequestDetailController()
	requestDetailGroup := router.Group("/request-details")
	{
		requestDetailGroup.GET("", requestDetailController.GetAllRequestDetails())
		requestDetailGroup.GET("/:id", requestDetailController.GetRequestDetailById())
		requestDetailGroup.POST("", requestDetailController.CreateRequestDetail())
		requestDetailGroup.PUT("", requestDetailController.UpdateRequestDetail())
		requestDetailGroup.DELETE("", requestDetailController.DeleteRequestDetail())
	}
}

func initInventoryAdjustmentRouter(router *gin.Engine) {
	adjustmentController := controller.NewInventoryAdjustmentController()
	adjustmentGroup := router.Group("/adjustments")
	{
		adjustmentGroup.GET("", adjustmentController.GetAllAdjustments())
		adjustmentGroup.GET("/:id", adjustmentController.GetAdjustmentById())
		adjustmentGroup.POST("", adjustmentController.CreateAdjustment())
		adjustmentGroup.PUT("", adjustmentController.UpdateAdjustment())
		adjustmentGroup.DELETE("", adjustmentController.DeleteAdjustment())
		adjustmentGroup.POST("/approval", adjustmentController.ApproveAdjustment())
	}
}

func initInventoryAuditRouter(router *gin.Engine) {
	auditController := controller.NewInventoryAuditController()
	auditGroup := router.Group("/audits")
	{
		auditGroup.GET("", auditController.GetAllAudits())
		auditGroup.GET("/:id", auditController.GetAuditById())
		auditGroup.POST("", auditController.CreateAudit())
		auditGroup.PUT("", auditController.UpdateAudit())
		auditGroup.DELETE("", auditController.DeleteAudit())
		auditGroup.POST("/approval", auditController.ApproveAudit())
		auditGroup.POST("/confirm", auditController.ConfirmAudit())
	}
}
func initInventoryAuditDetailRouter(router *gin.Engine) {
	auditDetailController := controller.NewInventoryAuditDetailController()
	auditDetailGroup := router.Group("/audit-details")
	{
		auditDetailGroup.GET("", auditDetailController.GetAllInventoryAuditDetails())
		// auditDetailGroup.GET("/:id", auditDetailController.GetInventoryAuditDetailById())
		auditDetailGroup.POST("", auditDetailController.CreateInventoryAuditDetail())
		auditDetailGroup.PUT("", auditDetailController.UpdateInventoryAuditDetail())
		auditDetailGroup.DELETE("", auditDetailController.DeleteInventoryAuditDetail())
	}
}

func initInventoryLedgerRouter(router *gin.Engine) {
	ledgerController := controller.NewInventoryLedgerController()
	ledgerGroup := router.Group("/ledgers")
	{
		ledgerGroup.GET("", ledgerController.GetAllLedgers())
		ledgerGroup.GET("/export", ledgerController.ExportLedgersExcel())
		ledgerGroup.GET("/:id", ledgerController.GetLedgerById())
	}
}

func initComponentCategoryRouter(router *gin.Engine) {
	categoryController := controller.NewComponentCategoryController()
	categoryGroup := router.Group("/categories")
	{
		categoryGroup.GET("", categoryController.GetAllCategories())
		categoryGroup.GET("/:id", categoryController.GetCategoryById())
		categoryGroup.POST("", categoryController.CreateCategory())
		categoryGroup.PUT("", categoryController.UpdateCategory())
		categoryGroup.DELETE("", categoryController.DeleteCategory())
	}
}

func initAuthRouter(router *gin.Engine) {
	authController := controller.NewAuthController()
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authController.GetToken())
		authGroup.POST("/refresh", authController.GetToken())
	}
}

func initUploadRouter(router *gin.Engine) {
	uploadController := controller.NewUploadController()
	uploadGroup := router.Group("/uploads")
	{
		uploadGroup.POST("/base64", uploadController.UploadBase64())
		uploadGroup.POST("/multipart", uploadController.UploadMultipart())
		uploadGroup.POST("/multiple", uploadController.UploadMultiple())
	}
}
