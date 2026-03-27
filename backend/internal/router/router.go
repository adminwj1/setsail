package router

import (
	"projecthub/internal/handler"
	"projecthub/internal/middleware"
	"projecthub/pkg/jwt"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "projecthub/docs"
)

type Router struct {
	engine  *gin.Engine
	jwtUtil *jwt.JWT

	authHandler        *handler.AuthHandler
	userHandler        *handler.UserHandler
	roleHandler        *handler.RoleHandler
	menuHandler        *handler.MenuHandler
	productHandler     *handler.ProductHandler
	requirementHandler *handler.RequirementHandler

	authMiddleware   *middleware.AuthMiddleware
	casbinMiddleware *middleware.CasbinMiddleware
}

func NewRouter(
	jwtUtil *jwt.JWT,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
	menuHandler *handler.MenuHandler,
	productHandler *handler.ProductHandler,
	requirementHandler *handler.RequirementHandler,
	authMiddleware *middleware.AuthMiddleware,
	casbinMiddleware *middleware.CasbinMiddleware,
) *Router {
	return &Router{
		engine:             gin.Default(),
		jwtUtil:            jwtUtil,
		authHandler:        authHandler,
		userHandler:        userHandler,
		roleHandler:        roleHandler,
		menuHandler:        menuHandler,
		productHandler:     productHandler,
		requirementHandler: requirementHandler,
		authMiddleware:     authMiddleware,
		casbinMiddleware:   casbinMiddleware,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(middleware.CORS())

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.engine.POST("/api/login", r.authHandler.Login)
	auth := r.engine.Group("/api/auth")
	{

		auth.POST("/logout", r.authHandler.Logout)
		auth.GET("/userinfo", r.authMiddleware.Handler(), r.authHandler.GetUserInfo)
	}

	api := r.engine.Group("/api")
	api.Use(r.authMiddleware.Handler())
	{
		menus := api.Group("/menus")
		menus.Use(r.casbinMiddleware.Handler())
		{
			menus.GET("/router", r.menuHandler.GetRouter)
			menus.GET("/tree", r.menuHandler.GetTree)
			menus.GET("", r.menuHandler.GetAll)
			menus.GET("/:id", r.menuHandler.GetByID)
			menus.POST("", r.menuHandler.Create)
			menus.PUT("/:id", r.menuHandler.Update)
			menus.DELETE("/:id", r.menuHandler.Delete)
		}

		roles := api.Group("/roles")
		roles.Use(r.casbinMiddleware.Handler())
		{
			roles.GET("", r.roleHandler.GetAll)
			roles.GET("/:id", r.roleHandler.GetByID)
			roles.POST("", r.roleHandler.Create)
			roles.PUT("/:id", r.roleHandler.Update)
			roles.DELETE("/:id", r.roleHandler.Delete)
			roles.GET("/:id/menus", r.roleHandler.GetMenus)
			roles.POST("/:id/menus", r.roleHandler.SetMenus)
		}

		users := api.Group("/users")
		users.Use(r.casbinMiddleware.Handler())
		{
			users.GET("", r.userHandler.GetAll)
			users.POST("", r.userHandler.Create)
			users.PUT("/:id", r.userHandler.Update)
			users.DELETE("/:id", r.userHandler.Delete)
		}

		products := api.Group("/products")
		products.Use(r.casbinMiddleware.Handler())
		{
			products.GET("", r.productHandler.GetAll)
			products.GET("/:id", r.productHandler.GetByID)
			products.POST("", r.productHandler.Create)
			products.PUT("/:id", r.productHandler.Update)
			products.DELETE("/:id", r.productHandler.Delete)
			products.GET("/:id/requirements", r.requirementHandler.GetByProductID)
		}

		requirements := api.Group("/requirements")
		requirements.Use(r.casbinMiddleware.Handler())
		{
			requirements.GET("", r.requirementHandler.GetAll)
			requirements.GET("/:id", r.requirementHandler.GetByID)
			requirements.POST("", r.requirementHandler.Create)
			requirements.PUT("/:id", r.requirementHandler.Update)
			requirements.DELETE("/:id", r.requirementHandler.Delete)
		}
	}

	return r.engine
}

// InitCasbin 初始化Casbin策略
func InitCasbin(e *casbin.Enforcer) error {
	// 添加默认策略：超级管理员拥有所有权限
	_, err := e.AddPolicy("1", "/api/menus", "GET")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/menus/:id", "GET")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/menus", "POST")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/menus/:id", "PUT")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/menus/:id", "DELETE")
	if err != nil {
		return err
	}

	_, err = e.AddPolicy("1", "/api/roles", "GET")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/roles/:id", "GET")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/roles", "POST")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/roles/:id", "PUT")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/roles/:id", "DELETE")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/roles/:id/menus", "GET")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/roles/:id/menus", "POST")
	if err != nil {
		return err
	}

	_, err = e.AddPolicy("1", "/api/users", "GET")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/users", "POST")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/users/:id", "PUT")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/users/:id", "DELETE")
	if err != nil {
		return err
	}

	_, err = e.AddPolicy("1", "/api/products", "GET")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/products/:id", "GET")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/products", "POST")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/products/:id", "PUT")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/products/:id", "DELETE")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/products/:id/requirements", "GET")
	if err != nil {
		return err
	}

	_, err = e.AddPolicy("1", "/api/requirements", "GET")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/requirements/:id", "GET")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/requirements", "POST")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/requirements/:id", "PUT")
	if err != nil {
		return err
	}
	_, err = e.AddPolicy("1", "/api/requirements/:id", "DELETE")
	if err != nil {
		return err
	}

	return nil
}
