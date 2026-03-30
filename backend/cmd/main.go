package main

import (
	"fmt"
	"log"
	"os"

	"projecthub/config"
	"projecthub/internal/handler"
	"projecthub/internal/middleware"
	"projecthub/internal/model"
	"projecthub/internal/repository"
	"projecthub/internal/router"
	"projecthub/internal/service"
	"projecthub/pkg/database"
	"projecthub/pkg/jwt"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()

	if _, err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	if err := initTables(); err != nil {
		log.Fatalf("Failed to initialize tables: %v", err)
	}

	if err := initDefaultData(); err != nil {
		log.Printf("Warning: Failed to initialize default data: %v", err)
	}

	jwtUtil := jwt.NewJWT(cfg.JWT.Secret, cfg.JWT.ExpireHour)

	// 使用 GORM 适配器创建 Casbin enforcer
	adapter, err := gormadapter.NewAdapterByDB(database.DB)
	if err != nil {
		log.Fatalf("Failed to create casbin adapter: %v", err)
	}

	dir, _ := os.Getwd()
	e, err := casbin.NewEnforcer(dir+"/"+"config/rbac_model.conf", adapter)
	if err != nil {
		log.Fatalf("Failed to create casbin enforcer: %v", err)
	}
	e.EnableAutoSave(true)

	if err := router.InitCasbin(e); err != nil {
		log.Printf("Warning: Failed to init casbin policy: %v", err)
	}

	casbinMiddleware := middleware.NewCasbinMiddleware(e)

	userRepo := repository.NewUserRepository()
	roleRepo := repository.NewRoleRepository()
	menuRepo := repository.NewMenuRepository()
	productRepo := repository.NewProductRepository()
	reqRepo := repository.NewRequirementRepository()

	authSvc := service.NewAuthService(userRepo, jwtUtil)
	userSvc := service.NewUserService(userRepo)
	roleSvc := service.NewRoleService(roleRepo, menuRepo, e)
	menuSvc := service.NewMenuService(menuRepo, roleRepo)
	productSvc := service.NewProductService(productRepo)
	reqSvc := service.NewRequirementService(reqRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	roleHandler := handler.NewRoleHandler(roleSvc)
	menuHandler := handler.NewMenuHandler(menuSvc)
	productHandler := handler.NewProductHandler(productSvc)
	reqHandler := handler.NewRequirementHandler(reqSvc)

	authMiddleware := middleware.NewAuthMiddleware(jwtUtil)

	r := router.NewRouter(
		jwtUtil,
		authHandler,
		userHandler,
		roleHandler,
		menuHandler,
		productHandler,
		reqHandler,
		authMiddleware,
		casbinMiddleware,
	)

	engine := r.Setup()

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initTables() error {
	return database.DB.AutoMigrate(
		&model.User{},
		&model.UserRole{},
		&model.Role{},
		&model.RoleMenu{},
		&model.Menu{},
		&model.Product{},
		&model.Requirement{},
	)
}

func initDefaultData() error {
	var count int64
	database.DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		log.Println("Default data already exists, skipping initialization")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin@123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: "admin",
		Password: string(hashedPassword),
		Nickname: "超级管理员",
		Email:    "admin@example.com",
		Status:   1,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return err
	}

	menus := []model.Menu{
		{ID: 1, ParentID: 0, Name: "系统管理", Code: "system", Path: "", ApiPath: "", Icon: "setting", Type: 1, Sort: 3, Status: 1},
		{ID: 2, ParentID: 1, Name: "用户管理", Code: "user:list", Path: "/system/user", ApiPath: "/api/users", Icon: "user", Type: 2, Sort: 1, Status: 1},
		{ID: 3, ParentID: 1, Name: "角色管理", Code: "role:list", Path: "/system/role", ApiPath: "/api/roles", Icon: "role", Type: 2, Sort: 2, Status: 1},
		{ID: 4, ParentID: 1, Name: "菜单管理", Code: "menu:list", Path: "/system/menu", ApiPath: "/api/menus", Icon: "menu", Type: 2, Sort: 3, Status: 1},
		{ID: 5, ParentID: 0, Name: "产品管理", Code: "product", Path: "", ApiPath: "", Icon: "product", Type: 1, Sort: 1, Status: 1},
		{ID: 6, ParentID: 5, Name: "产品列表", Code: "product:list", Path: "/product/list", ApiPath: "/api/products", Icon: "list", Type: 2, Sort: 1, Status: 1},
		{ID: 7, ParentID: 0, Name: "需求管理", Code: "requirement", Path: "", ApiPath: "", Icon: "requirement", Type: 1, Sort: 2, Status: 1},
		{ID: 8, ParentID: 7, Name: "需求列表", Code: "requirement:list", Path: "/requirement/list", ApiPath: "/api/requirements", Icon: "list", Type: 2, Sort: 1, Status: 1},
	}

	for _, m := range menus {
		if err := database.DB.Create(&m).Error; err != nil {
			fmt.Printf("err: %v\n", err)
			return err
		}
	}

	role := &model.Role{
		Name:   "超级管理员",
		Code:   "super_admin",
		Desc:   "拥有所有权限",
		Sort:   1,
		Status: 1,
	}

	if err := database.DB.Create(role).Error; err != nil {
		return err
	}

	if err := database.DB.Create(&model.UserRole{UserID: 1, RoleID: role.ID}).Error; err != nil {
		return err
	}

	menuIDs := []uint{1, 2, 3, 4, 5, 6, 7, 8}
	for _, menuID := range menuIDs {
		if err := database.DB.Create(&model.RoleMenu{RoleID: role.ID, MenuID: menuID}).Error; err != nil {
			return err
		}
	}

	return nil
}
