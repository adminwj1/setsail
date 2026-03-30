package router

// Route 定义路由配置
type Route struct {
	Path   string
	Method string
}

// Routes 定义所有需要权限控制的路由
var Routes = []Route{
	// 菜单路由
	{Path: "/api/menus", Method: "GET"},
	{Path: "/api/menus/router", Method: "GET"},
	{Path: "/api/menus/tree", Method: "GET"},
	{Path: "/api/menus/:id", Method: "GET"},
	{Path: "/api/menus", Method: "POST"},
	{Path: "/api/menus/:id", Method: "PUT"},
	{Path: "/api/menus/:id", Method: "DELETE"},

	// 角色路由
	{Path: "/api/roles", Method: "GET"},
	{Path: "/api/roles/:id", Method: "GET"},
	{Path: "/api/roles", Method: "POST"},
	{Path: "/api/roles/:id", Method: "PUT"},
	{Path: "/api/roles/:id", Method: "DELETE"},
	{Path: "/api/roles/:id/menus", Method: "GET"},
	{Path: "/api/roles/:id/menus", Method: "POST"},

	// 用户路由
	{Path: "/api/users", Method: "GET"},
	{Path: "/api/users", Method: "POST"},
	{Path: "/api/users/:id", Method: "PUT"},
	{Path: "/api/users/:id", Method: "DELETE"},

	// 产品路由
	{Path: "/api/products", Method: "GET"},
	{Path: "/api/products/:id", Method: "GET"},
	{Path: "/api/products", Method: "POST"},
	{Path: "/api/products/:id", Method: "PUT"},
	{Path: "/api/products/:id", Method: "DELETE"},
	{Path: "/api/products/:id/requirements", Method: "GET"},

	// 需求路由
	{Path: "/api/requirements", Method: "GET"},
	{Path: "/api/requirements/:id", Method: "GET"},
	{Path: "/api/requirements", Method: "POST"},
	{Path: "/api/requirements/:id", Method: "PUT"},
	{Path: "/api/requirements/:id", Method: "DELETE"},
}
