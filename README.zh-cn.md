# RMS Survey

<p align="center">
  <strong>AI 驱动的开源问卷调查与考试系统</strong>
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
- **在线考试** - 支持自动评分、练习记录、错题本
- **AI 聊天集成** - 可选的 AI 聊天流接口（需配置）
- **审批流程** - 简单的答卷审批状态跟踪
- **模板库** - 可复用的题目模板，支持分类与标签
- **题库管理** - 支持从 Excel 批量导入导出题目
- **数据看板** - 实时统计与数据可视化
- **附件上传** - 题目支持上传附件

## 技术栈

### 后端 (server/)
- **Go 1.24** + Gin Web 框架
- **GORM** MySQL ORM
- **JWT** 认证 + RSA 密钥加密
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
│   │   ├── pkg/               # 工具库（jwt、rsa、cache 等）
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
        ├── router/            # Vue Router
        ├── stores/            # Pinia 状态管理
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

ai:
  api_key: ""
  base_url: ""
  model: ""
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
| `AI_API_KEY` | AI 服务 API 密钥（可选） |

## 配置说明

### 服务配置 (config.yaml)

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `server.port` | HTTP 服务端口 | 8080 |
| `database.dsn` | MySQL 连接字符串 | 必填 |
| `jwt.secret` | JWT 签名密钥 | 必填 |
| `jwt.cookie_name` | 认证 Cookie 名称 | survey_token |
| `storage.local_path` | 文件上传目录 | ./uploads |
| `ai.api_key` | AI 服务 API 密钥 | 可选 |
| `ai.base_url` | AI 服务地址 | 可选 |
| `ai.model` | AI 模型名称 | 可选 |

## API 参考

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/public/rsaPublicKey` | 获取 RSA 公钥（用于密码加密） |
| `POST` | `/api/public/login` | 用户登录 |
| `POST` | `/api/public/register` | 用户注册 |
| `POST` | `/api/public/loadProject` | 加载问卷/考试 |
| `POST` | `/api/public/saveAnswer` | 提交答卷 |
| `POST` | `/api/public/tempSaveAnswer` | 临时保存答卷 |
| `POST` | `/api/public/uploadAttachment` | 上传题目附件 |
| `GET` | `/api/public/preview/:attachmentId` | 预览附件 |

### 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET/POST` | `/api/project/*` | 项目增删改查 |
| `GET/POST` | `/api/survey/*` | 问卷设置与逻辑 |
| `GET/POST` | `/api/answer/*` | 答卷管理 |
| `GET/POST` | `/api/template/*` | 模板增删改查（含分类标签） |
| `GET/POST` | `/api/repo/*` | 题库管理 |
| `GET/POST` | `/api/repo/userBook/*` | 用户错题本/收藏夹 |
| `GET/POST` | `/api/repo/import` | 从 Excel 导入题目 |
| `GET` | `/api/repo/export` | 导出题目到 Excel |
| `GET/POST` | `/api/workflow/*` | 审批流程操作 |
| `GET` | `/api/ai/chat/models` | 获取可用 AI 模型 |
| `GET` | `/api/ai/chat/stream` | AI 聊天流（SSE） |
| `GET/POST` | `/api/file/*` | 文件上传下载 |
| `GET/POST` | `/api/dashboard/*` | 看板数据 |
| `GET/POST` | `/api/exercise/*` | 练习记录 |
| `GET` | `/api/report/:shortId` | 报表数据 |
| `GET/POST` | `/api/system/*` | 系统设置 |

## 数据库

数据表使用 `t_` 前缀。核心表：

| 表名 | 说明 |
|------|------|
| `t_project` | 问卷和考试 |
| `t_answer` | 答卷记录 |
| `t_template` | 题目模板 |
| `t_repo` | 题库 |
| `t_user_book` | 用户题目收藏（错题本） |
| `t_flow_operation` | 审批流程历史 |
| `t_file` | 上传的文件 |
| `t_user` | 用户 |
| `t_role` | 角色 |
| `t_account` | 认证账户 |
| `t_dashboard` | 看板配置 |
| `t_sys_info` | 系统设置 |

## 许可证

[MIT License](LICENSE)

