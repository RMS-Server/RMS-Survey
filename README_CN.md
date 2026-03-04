# RMS Survey

<p align="center">
  <strong>开源问卷调查系统</strong>
</p>

<p align="center">
  <a href="#功能特性">功能特性</a> •
  <a href="#技术栈">技术栈</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#配置说明">配置说明</a> •
  <a href="#api-参考">API 参考</a> •
  <a href="#许可证">许可证</a>
</p>

---

## 功能特性

- **问卷设计器** - 拖拽式问卷设计，支持多种题型
- **模板库** - 可复用的题目模板，快速创建问卷
- **参与者管理** - 控制问卷的访问权限
- **答卷管理** - 查看、删除、恢复、导出答卷，支持已读/未读状态追踪
- **问卷逻辑** - 条件显示和跳题逻辑，支持 AND/OR 组合条件
- **附件上传** - 题目支持上传附件，图片支持内联预览
- **IP限制** - 支持按IP限制提交次数和提交间隔
- **计时追踪** - 记录每道题的答题时长和总耗时
- **草稿保存** - 自动保存答题进度到本地存储，页面刷新后恢复
- **回收站** - 项目和答卷软删除，支持恢复
- **角色权限** - 用户管理与角色控制，支持 OAuth 2.0 PKCE 单点登录
- **玻璃态界面** - 现代玻璃效果设计，支持鼠标跟随发光效果

## 技术栈

### 后端 (server/)
- **Go 1.24** + Gin Web 框架
- **GORM** MySQL ORM
- **JWT** 认证 + OAuth 2.0 PKCE 单点登录
- **Viper** 配置管理

### 前端 (website/)
- **Vue 3** Composition API
- **TypeScript** 类型安全
- **Vite** 快速开发构建
- **Pinia** 状态管理
- **Ant Design Vue** UI 组件库

## 项目结构

```
rms-survey/
├── server/                    # Go 后端
│   ├── cmd/server/main.go     # 入口文件
│   ├── configs/               # 配置文件
│   ├── internal/
│   │   ├── config/            # 配置加载
│   │   ├── dto/               # 请求响应结构体
│   │   ├── handler/           # HTTP 处理器
│   │   ├── middleware/        # 中间件（认证、跨域等）
│   │   ├── model/             # GORM 模型
│   │   ├── pkg/               # 工具库（jwt、cache 等）
│   │   ├── repository/        # 数据访问层
│   │   ├── router/            # 路由注册
│   │   └── service/           # 业务逻辑层
│   ├── static/                # SPA 静态文件
│   └── uploads/               # 文件存储
│
└── website/                   # Vue 3 前端
    └── src/
        ├── api/               # Axios API 客户端
        ├── components/        # Vue 组件
        ├── composables/       # Vue 组合式函数（useLogicEvaluator、useTimingTracker、useDraftManager）
        ├── router/            # Vue Router
        ├── stores/            # Pinia 状态管理
        ├── styles/            # CSS 样式（theme.css、glassmorphism.css）
        ├── types/             # TypeScript 类型定义
        └── views/             # 页面组件
```

## 快速开始

### 环境要求

- Go 1.24+
- Node.js 18+
- MySQL 8.0+

### 后端启动

```bash
cd server

# 安装依赖
go mod download

# 创建配置文件
mkdir -p configs
cat > configs/config.yaml << EOF
server:
  port: 8080

database:
  dsn: "用户名:密码@tcp(localhost:3306)/survey_db?charset=utf8mb4&parseTime=True&loc=Local"

jwt:
  secret: "至少64字符的JWT密钥以确保安全"
  cookie_name: "survey_token"

storage:
  local_path: "./uploads"
EOF

# 创建上传目录
mkdir -p uploads

# 启动服务
go run cmd/server/main.go
```

### 前端启动

```bash
cd website

# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
npm run build

# 类型检查
npm run type-check
```

### 环境变量

| 变量 | 说明 |
|------|------|
| `JWT_SECRET` | JWT 签名密钥（至少 64 字符） |

## 配置说明

### 服务配置 (config.yaml)

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `server.port` | HTTP 服务端口 | 8080 |
| `database.dsn` | MySQL 连接字符串 | 必填 |
| `jwt.secret` | JWT 签名密钥 | 必填 |
| `jwt.cookie_name` | 认证 Cookie 名称 | survey_token |
| `storage.local_path` | 文件上传目录 | ./uploads |

### OAuth 2.0 配置

| 字段 | 说明 |
|------|------|
| `oauth.enabled` | 启用 OAuth 认证 |
| `oauth.client_id` | SSO 提供商分配的客户端 ID |
| `oauth.auth_url` | 授权端点 URL |
| `oauth.token_url` | 令牌端点 URL |
| `oauth.userinfo_url` | 用户信息端点 URL |
| `oauth.redirect_url` | 应用回调 URL |
| `oauth.scopes` | OAuth 权限范围（空格分隔） |
| `oauth.min_permission_level` | 最低权限等级要求 |

## API 参考

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/public/login` | 用户登录 |
| `POST` | `/api/public/logout` | 用户登出 |
| `GET` | `/api/oauth/authorize` | 获取 OAuth 授权配置 |
| `GET` | `/api/oauth/callback` | OAuth 回调（返回 HTML 用于 PKCE 流程） |
| `POST` | `/api/oauth/callback` | 完成 OAuth 登录并获取 JWT |
| `POST` | `/api/oauth/refresh` | 刷新 JWT 令牌 |
| `POST` | `/api/oauth/logout` | OAuth 登出（清除会话） |
| `GET` | `/api/public/listRegisterRole` | 获取可注册角色列表 |
| `POST` | `/api/public/loadProject` | 加载问卷 |
| `POST` | `/api/public/validateProject` | 验证问卷访问权限 |
| `POST` | `/api/public/statistics` | 获取问卷统计 |
| `POST` | `/api/public/saveAnswer` | 提交答卷 |
| `POST` | `/api/public/tempSaveAnswer` | 临时保存答卷 |
| `POST` | `/api/public/uploadAttachment` | 上传题目附件 |
| `GET` | `/api/public/preview/:attachmentId` | 预览附件 |
| `GET` | `/captcha/get` | 获取验证码图片 |
| `POST` | `/captcha/check` | 验证验证码 |

### 认证接口

#### 用户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/currentUser` | 获取当前用户信息 |
| `GET` | `/userOverview` | 获取用户概览统计 |
| `GET` | `/api/user/list` | 用户列表（分页） |
| `POST` | `/api/user` | 创建用户 |
| `PUT` | `/api/user` | 更新用户 |
| `GET` | `/api/user/:id` | 获取用户详情 |
| `DELETE` | `/api/user/:id` | 删除用户 |
| `POST` | `/api/user/bindRole` | 绑定用户角色 |
| `PUT` | `/api/user/updatePassword` | 修改密码 |
| `POST` | `/importUser` | 从 Excel 导入用户 |

#### 项目管理

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/project/list` | 项目列表（分页） |
| `GET` | `/api/project` | 获取项目详情 |
| `GET` | `/api/project/setting` | 获取项目设置 |
| `POST` | `/api/project/create` | 创建项目 |
| `POST` | `/api/project/update` | 更新项目 |
| `POST` | `/api/project/delete` | 删除项目（软删除） |
| `GET` | `/api/project/trash` | 已删除项目列表 |
| `POST` | `/api/project/destroy` | 彻底删除项目 |
| `POST` | `/api/project/restore` | 恢复已删除项目 |

#### 参与者管理

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/project/partner/list` | 参与者列表 |
| `POST` | `/api/project/partner/create` | 添加参与者 |
| `POST` | `/api/project/partner/delete` | 移除参与者 |
| `GET` | `/api/project/partner/download` | 下载参与者 Excel |
| `POST` | `/api/project/partner/import` | 导入参与者 Excel |

#### 选择器

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/project/selectUser` | 用户选择器 |
| `POST` | `/api/project/selectRole` | 角色选择器 |
| `POST` | `/api/project/selectTemplate` | 模板选择器 |

#### 问卷设置

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET/POST` | `/api/survey/setting` | 获取/更新问卷设置 |
| `GET/POST` | `/api/survey/logic` | 获取/更新问卷逻辑规则 |

#### 问卷逻辑规则

逻辑规则实现题目条件显示功能。每条规则包含：
- **目标题目**：根据条件显示/隐藏的题目
- **条件**：一个或多个条件，支持等于、包含、大于、小于等操作符
- **逻辑运算符**：AND（所有条件都满足）或 OR（任一条件满足）

**前端组件：**
- `LogicRuleEditor.vue` — 逻辑规则编辑弹窗
- `ConditionBuilder.vue` — 可视化条件构建器，支持 AND/OR 组合
- `useLogicEvaluator.ts` — 运行时逻辑评估组合式函数

规则存储在项目设置中，在问卷填写时由前端实时评估。隐藏的题目在提交时跳过验证。当引用的题目被删除时，相关规则会自动清理。

#### 答卷管理

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/answer/list` | 答卷列表（分页，支持跨项目查询） |
| `GET` | `/api/answer/trash` | 已删除答卷列表 |
| `GET` | `/api/answer` | 获取答卷详情 |
| `POST` | `/api/answer/create` | 创建答卷 |
| `POST` | `/api/answer/update` | 更新答卷 |
| `POST` | `/api/answer/delete` | 删除答卷（软删除） |
| `POST` | `/api/answer/destroy` | 彻底删除答卷 |
| `POST` | `/api/answer/restore` | 恢复已删除答卷 |
| `POST` | `/api/answer/read` | 标记答卷为已读 |
| `POST` | `/api/answer/unread` | 标记答卷为未读 |
| `GET` | `/api/answer/download` | 导出答卷 Excel（含IP和计时数据） |

#### 模板管理

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/template/list` | 模板列表 |
| `GET` | `/api/template` | 获取模板详情 |
| `POST` | `/api/template/create` | 创建模板 |
| `POST` | `/api/template/update` | 更新模板 |
| `POST` | `/api/template/delete` | 删除模板 |
| `GET` | `/api/template/category/list` | 模板分类列表 |
| `GET` | `/api/template/tag/list` | 模板标签列表 |

#### 系统管理

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/system` | 获取系统信息 |
| `POST` | `/api/system/update` | 更新系统设置 |
| `GET` | `/api/system/role/list` | 角色列表 |
| `POST` | `/api/system/role/create` | 创建角色 |
| `POST` | `/api/system/role/update` | 更新角色 |
| `POST` | `/api/system/role/delete` | 删除角色 |

## 数据库

数据表使用 `t_` 前缀。核心表：

| 表名 | 说明 |
|------|------|
| `t_project` | 问卷项目 |
| `t_answer` | 答卷记录 |
| `t_template` | 题目模板 |
| `t_file` | 上传的文件 |
| `t_user` | 用户 |
| `t_role` | 角色 |
| `t_account` | 认证账户 |
| `t_oauth_session` | OAuth 刷新令牌 |
| `t_project_partner` | 问卷参与者 |
| `t_user_role` | 用户角色关联 |
| `t_sys_info` | 系统设置 |
| `t_comm_dict_item` | 字典项（问卷下拉选项） |

### 答卷模型字段

`t_answer` 表包含以下重要字段：
- `is_read`、`read_at`、`read_by` — 已读状态追踪
- `ip_address` — 提交者IP，用于限制校验
- `timing_info` — JSON字段，存储每题计时数据

## 许可证

[MIT License](LICENSE)
