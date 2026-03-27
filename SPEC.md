# 项目管理工具 - 规格说明书

## 1. 项目概述

**项目名称**: ProjectHub
**项目类型**: 前后端分离的项目管理平台
**核心功能**: 企业级项目管理，支持产品管理、需求管理、RBAC权限控制
**技术栈**:
- 后端: Golang + Gin + MySQL 8.0 + Casbin
- 前端: Vue3 + Element Plus
- 认证: JWT

---

## 2. 功能模块

### 2.1 用户认证模块

| 功能 | 说明 |
|------|------|
| 用户登录 | POST /api/auth/login，参数: username, password |
| 用户登出 | POST /api/auth/logout |
| 刷新Token | POST /api/auth/refresh |
| 获取当前用户信息 | GET /api/auth/userinfo |

**初始账号**:
- 超级管理员: username=admin, password=admin@123

### 2.2 权限管理模块 (RBAC)

使用 casbin 实现基于角色的访问控制。

**数据库表**:
- `sys_user` - 用户表
- `sys_role` - 角色表
- `sys_menu` - 菜单/权限表
- `sys_user_role` - 用户-角色关联表
- `sys_role_menu` - 角色-菜单关联表

**菜单类型**:
- `directory` - 目录
- `menu` - 菜单
- `button` - 按钮权限

**API接口**:
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/menus/router | 获取当前用户的动态路由菜单 |
| GET | /api/roles | 获取角色列表 |
| POST | /api/roles | 创建角色 |
| PUT | /api/roles/:id | 更新角色 |
| DELETE | /api/roles/:id | 删除角色 |
| GET | /api/menus | 获取所有菜单树 |
| POST | /api/menus | 创建菜单 |
| PUT | /api/menus/:id | 更新菜单 |
| DELETE | /api/menus/:id | 删除菜单 |
| GET | /api/users | 获取用户列表 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/:id | 更新用户 |
| DELETE | /api/users/:id | 删除用户 |

### 2.3 产品管理模块

**数据库表**:
- `pm_product` - 产品表

**产品表字段**:
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| name | varchar(100) | 产品名称 |
| code | varchar(50) | 产品代号 |
| description | text | 产品描述 |
| status | tinyint | 状态: 0-规划中, 1-开发中, 2-上线, 3-下线 |
| owner_id | bigint | 负责人ID |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

**API接口**:
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/products | 获取产品列表 |
| POST | /api/products | 创建产品 |
| GET | /api/products/:id | 获取产品详情 |
| PUT | /api/products/:id | 更新产品 |
| DELETE | /api/products/:id | 删除产品 |

### 2.4 需求管理模块

**数据库表**:
- `pm_requirement` - 需求表

**需求表字段**:
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| product_id | bigint | 所属产品ID |
| title | varchar(200) | 需求标题 |
| type | tinyint | 类型: 1-业务需求, 2-用户需求, 3-研发需求 |
| priority | tinyint | 优先级: 1-高, 2-中, 3-低 |
| status | tinyint | 状态: 0-待评审, 1-已采纳, 2-开发中, 3-已完成, 4-已拒绝 |
| description | text | 需求描述 |
| creator_id | bigint | 创建人ID |
| assignee_id | bigint | 负责人ID |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

**API接口**:
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/requirements | 获取需求列表 |
| POST | /api/requirements | 创建需求 |
| GET | /api/requirements/:id | 获取需求详情 |
| PUT | /api/requirements/:id | 更新需求 |
| DELETE | /api/requirements/:id | 删除需求 |
| GET | /api/products/:id/requirements | 获取指定产品的需求列表 |

---

## 3. 数据库设计

### 3.1 ER图概览

```
sys_user (1) ----< (N) sys_user_role >---- (N) sys_role (1)
                                                    |
                                                    | (N)
                                              sys_role_menu >---- (N) sys_menu

pm_product (1) ----< (N) pm_requirement
```

### 3.2 初始化数据

**超级管理员角色**:
- 角色名: super_admin
- 权限: 全部菜单权限

**默认菜单结构**:
```
├── 系统管理 (directory)
│   ├── 用户管理 (menu)
│   ├── 角色管理 (menu)
│   └── 菜单管理 (menu)
├── 产品管理 (directory)
│   └── 产品列表 (menu)
└── 需求管理 (directory)
    └── 需求列表 (menu)
```

---

## 4. API响应格式

**成功响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

**错误响应**:
```json
{
  "code": 400,
  "message": "error message",
  "data": null
}
```

**分页响应**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "pageSize": 10
  }
}
```

---

## 5. 认证机制

- 使用 JWT (JSON Web Token)
- Token 有效期: 24小时
-  Refresh Token 有效期: 7天
-  Token 放在 HTTP Header: `Authorization: Bearer <token>`

---

## 6. 验收标准

### 6.1 登录功能
- [ ] 管理员可以使用 admin/admin@123 登录
- [ ] 登录成功后返回 JWT token
- [ ] 错误的用户名或密码返回错误提示

### 6.2 权限功能
- [ ] 不同角色看到不同的菜单
- [ ] 动态路由根据后端接口渲染
- [ ] 无权限访问接口返回 403

### 6.3 产品管理
- [ ] 可以创建、查看、编辑、删除产品
- [ ] 产品列表支持分页

### 6.4 需求管理
- [ ] 可以创建、查看、编辑、删除需求
- [ ] 可以按产品筛选需求
- [ ] 需求支持类型、优先级、状态筛选

---

## 7. 项目结构

```
dome/
├── backend/                 # 后端项目
│   ├── cmd/                 # main.go入口
│   ├── config/              # 配置
│   ├── internal/             # 内部包
│   │   ├── handler/         # HTTP处理层
│   │   ├── middleware/       # 中间件
│   │   ├── model/           # 数据模型
│   │   ├── repository/      # 数据访问层
│   │   ├── service/         # 业务逻辑层
│   │   └── router/          # 路由
│   ├── pkg/                  # 公共包
│   │   ├── database/        # 数据库连接
│   │   ├── jwt/             # JWT工具
│   │   └── response/        # 统一响应
│   ├── go.mod
│   └── go.sum
├── frontend/                 # 前端项目 (可选)
└── SPEC.md
```
