# RMS Survey

<p align="center">
  <strong>Open Source Survey System</strong>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#tech-stack">Tech Stack</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#api-reference">API Reference</a> •
  <a href="#license">License</a>
</p>

---

## Features

- **Survey Builder** - Drag-and-drop questionnaire designer with multiple question types
- **Template Library** - Reusable question templates with categories and tags
- **Participant Management** - Control who can access your surveys
- **Response Management** - View, delete, restore, and export survey responses
- **Survey Logic** - Conditional question display with AND/OR logic rules
- **Question Attachments** - Upload files attached to questions with inline preview
- **Recycle Bin** - Soft delete with restore capability for projects and answers
- **Role-based Access** - User management with role permissions and OAuth 2.0 PKCE SSO authentication

## Tech Stack

### Backend (server/)
- **Go 1.24** with Gin web framework
- **GORM** for MySQL database ORM
- **JWT** authentication with OAuth 2.0 PKCE SSO integration
- **Viper** for configuration management

### Frontend (website/)
- **Vue 3** with Composition API and `<script setup>` syntax
- **TypeScript** for type safety
- **Vite** for fast development and building
- **Pinia** for state management
- **Ant Design Vue** for UI components

## Project Structure

```
rms-survey/
├── server/                    # Go backend
│   ├── cmd/server/main.go     # Entry point
│   ├── configs/               # Configuration files
│   ├── internal/
│   │   ├── config/            # Config loading
│   │   ├── dto/               # Request/response DTOs
│   │   ├── handler/           # HTTP handlers
│   │   ├── middleware/        # Auth, CORS, etc.
│   │   ├── model/             # GORM models
│   │   ├── pkg/               # Utilities (jwt, cache, response, storage)
│   │   ├── repository/        # Data access layer
│   │   ├── router/            # Route registration
│   │   └── service/           # Business logic
│   ├── static/                # SPA static files
│   └── uploads/               # File storage
│
└── website/                   # Vue 3 frontend
    └── src/
        ├── api/               # Axios API clients
        ├── components/        # Vue components
        ├── router/            # Vue Router
        ├── stores/            # Pinia stores
        ├── types/             # TypeScript interfaces
        └── views/             # Page components
```

## Quick Start

### Prerequisites

- Go 1.24+
- Node.js 18+
- MySQL 8.0+

### Backend Setup

```bash
cd server

# Install dependencies
go mod download

# Create config file
mkdir -p configs
cat > configs/config.yaml << EOF
server:
  port: 8080

database:
  dsn: "your_user:your_password@tcp(localhost:3306)/survey_db?charset=utf8mb4&parseTime=True&loc=Local"

jwt:
  secret: "your-jwt-secret-at-least-64-characters-long-for-security"
  cookie_name: "survey_token"

storage:
  local_path: "./uploads"
EOF

# Create upload directory
mkdir -p uploads

# Run the server
go run cmd/server/main.go
```

### Frontend Setup

```bash
cd website

# Install dependencies
npm install

# Development server
npm run dev

# Production build
npm run build

# Type check
npm run type-check
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `JWT_SECRET` | JWT signing secret (min 64 chars) |

## Configuration

### Server Config (config.yaml)

| Field | Description | Default |
|-------|-------------|---------|
| `server.port` | HTTP server port | 8080 |
| `database.dsn` | MySQL connection string | Required |
| `jwt.secret` | JWT signing secret | Required |
| `jwt.cookie_name` | Auth cookie name | survey_token |
| `storage.local_path` | File upload directory | ./uploads |

### OAuth 2.0 Configuration

| Field | Description |
|-------|-------------|
| `oauth.enabled` | Enable OAuth authentication |
| `oauth.client_id` | OAuth client ID from SSO provider |
| `oauth.auth_url` | Authorization endpoint URL |
| `oauth.token_url` | Token endpoint URL |
| `oauth.userinfo_url` | User info endpoint URL |
| `oauth.redirect_url` | Callback URL for your application |
| `oauth.scopes` | OAuth scopes (space-separated) |
| `oauth.min_permission_level` | Minimum permission level required |

## API Reference

### Public Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/public/login` | User login |
| `POST` | `/api/public/logout` | User logout |
| `GET` | `/api/oauth/authorize` | Get OAuth authorization config |
| `GET` | `/api/oauth/callback` | OAuth callback (returns HTML for PKCE flow) |
| `POST` | `/api/oauth/callback` | Complete OAuth login and get JWT |
| `POST` | `/api/oauth/refresh` | Refresh JWT token |
| `POST` | `/api/oauth/logout` | OAuth logout (clear session) |
| `GET` | `/api/public/listRegisterRole` | Get available registration roles |
| `POST` | `/api/public/loadProject` | Load survey by code |
| `POST` | `/api/public/validateProject` | Validate survey access |
| `POST` | `/api/public/statistics` | Get survey statistics |
| `POST` | `/api/public/saveAnswer` | Submit survey answer |
| `POST` | `/api/public/tempSaveAnswer` | Temporarily save answer |
| `POST` | `/api/public/uploadAttachment` | Upload question attachment |
| `GET` | `/api/public/preview/:attachmentId` | Preview attachment file |
| `GET` | `/captcha/get` | Get captcha image |
| `POST` | `/captcha/check` | Verify captcha |

### Protected Endpoints

#### User Management

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/currentUser` | Get current user info |
| `GET` | `/userOverview` | Get user overview statistics |
| `GET` | `/api/user/list` | List users (paginated) |
| `POST` | `/api/user` | Create user |
| `PUT` | `/api/user` | Update user |
| `GET` | `/api/user/:id` | Get user by ID |
| `DELETE` | `/api/user/:id` | Delete user |
| `POST` | `/api/user/bindRole` | Bind user role |
| `PUT` | `/api/user/updatePassword` | Update password |
| `POST` | `/importUser` | Import users from Excel |

#### Project Management

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/project/list` | List projects (paginated) |
| `GET` | `/api/project` | Get project by ID |
| `GET` | `/api/project/setting` | Get project settings |
| `POST` | `/api/project/create` | Create project |
| `POST` | `/api/project/update` | Update project |
| `POST` | `/api/project/delete` | Delete project (soft) |
| `GET` | `/api/project/trash` | List deleted projects |
| `POST` | `/api/project/destroy` | Permanently delete project |
| `POST` | `/api/project/restore` | Restore deleted project |

#### Participant Management

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/project/partner/list` | List participants |
| `POST` | `/api/project/partner/create` | Add participant |
| `POST` | `/api/project/partner/delete` | Remove participant |
| `GET` | `/api/project/partner/download` | Download participants Excel |
| `POST` | `/api/project/partner/import` | Import participants Excel |

#### Selectors

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/project/selectUser` | User selector |
| `POST` | `/api/project/selectRole` | Role selector |
| `POST` | `/api/project/selectTemplate` | Template selector |

#### Survey Settings

| Method | Path | Description |
|--------|------|-------------|
| `GET/POST` | `/api/survey/setting` | Get/update survey settings |
| `GET/POST` | `/api/survey/logic` | Get/update survey logic rules |

#### Survey Logic Rules

Logic rules enable conditional question display. Each rule consists of:
- **Target question**: The question to show/hide based on conditions
- **Conditions**: One or more conditions with operator (equals, contains, greater than, less than, etc.)
- **Logic operator**: AND (all conditions must match) or OR (any condition matches)

**Frontend Components:**
- `LogicRuleEditor.vue` — Modal editor for creating/editing logic rules
- `ConditionBuilder.vue` — Visual condition builder for AND/OR logic
- `useLogicEvaluator.ts` — Runtime composable that evaluates rules during survey fill

Rules are stored in project settings and evaluated client-side. Hidden questions are excluded from validation during submission. Rules are automatically cleaned up when referenced questions are deleted.

#### Answer Management

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/answer/list` | List answers (paginated) |
| `GET` | `/api/answer/trash` | List deleted answers |
| `GET` | `/api/answer` | Get answer by ID |
| `POST` | `/api/answer/create` | Create answer |
| `POST` | `/api/answer/update` | Update answer |
| `POST` | `/api/answer/delete` | Delete answer (soft) |
| `POST` | `/api/answer/destroy` | Permanently delete answer |
| `POST` | `/api/answer/restore` | Restore deleted answer |
| `GET` | `/api/answer/download` | Export answers to Excel |

#### Template Management

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/template/list` | List templates |
| `GET` | `/api/template` | Get template by ID |
| `POST` | `/api/template/create` | Create template |
| `POST` | `/api/template/update` | Update template |
| `POST` | `/api/template/delete` | Delete template |
| `GET` | `/api/template/category/list` | List template categories |
| `GET` | `/api/template/tag/list` | List template tags |

#### System Management

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/system` | Get system info |
| `POST` | `/api/system/update` | Update system settings |
| `GET` | `/api/system/role/list` | List roles |
| `POST` | `/api/system/role/create` | Create role |
| `POST` | `/api/system/role/update` | Update role |
| `POST` | `/api/system/role/delete` | Delete role |

#### File Management

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/public/uploadAttachment` | Upload question attachment |
| `GET` | `/api/public/preview/:attachmentId` | Preview/download attachment |

## Database

Tables use `t_` prefix. Core tables:

| Table | Description |
|-------|-------------|
| `t_project` | Survey projects |
| `t_answer` | Survey responses |
| `t_template` | Question templates |
| `t_file` | Uploaded files |
| `t_user` | Users |
| `t_role` | Roles |
| `t_account` | Authentication accounts |
| `t_oauth_session` | OAuth refresh tokens |
| `t_project_partner` | Survey participants |
| `t_user_role` | User-role associations |
| `t_sys_info` | System settings |
| `t_comm_dict_item` | Dictionary items for dropdowns |

## License

[MIT License](LICENSE)
