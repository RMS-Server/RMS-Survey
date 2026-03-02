# RMS Survey

<p align="center">
  <strong>AI-Powered Open Source Survey & Exam System</strong>
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
- **Exam System** - Online examination with auto-scoring, exercise history and practice mode
- **AI Chat Integration** - Optional AI chat stream integration (configurable)
- **Approval Workflow** - Simple approval flow tracking for survey responses
- **Template Library** - Reusable question templates with categories and tags
- **Question Bank** - Centralized repository with Excel import/export support
- **Dashboard & Reports** - Real-time analytics and data visualization
- **File Attachments** - Upload attachments for survey questions

## Tech Stack

### Backend (server/)
- **Go 1.24** with Gin web framework
- **GORM** for MySQL database ORM
- **JWT** authentication with RSA key encryption
- **Viper** for configuration management

### Frontend (website/)
- **Vue 3** with Composition API
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
│   │   ├── pkg/               # Utilities (jwt, rsa, cache, etc.)
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

ai:
  api_key: ""
  base_url: ""
  model: ""
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
| `AI_API_KEY` | AI service API key (optional) |

## Configuration

### Server Config (config.yaml)

| Field | Description | Default |
|-------|-------------|---------|
| `server.port` | HTTP server port | 8080 |
| `database.dsn` | MySQL connection string | Required |
| `jwt.secret` | JWT signing secret | Required |
| `jwt.cookie_name` | Auth cookie name | survey_token |
| `storage.local_path` | File upload directory | ./uploads |
| `ai.api_key` | AI service API key | Optional |
| `ai.base_url` | AI service base URL | Optional |
| `ai.model` | AI model name | Optional |

## API Reference

### Public Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/public/rsaPublicKey` | Get RSA public key for password encryption |
| `POST` | `/api/public/login` | User login |
| `POST` | `/api/public/register` | User registration |
| `POST` | `/api/public/loadProject` | Load survey/exam by ID |
| `POST` | `/api/public/saveAnswer` | Submit survey/exam answer |
| `POST` | `/api/public/tempSaveAnswer` | Temporarily save answer |
| `POST` | `/api/public/uploadAttachment` | Upload question attachment |
| `GET` | `/api/public/preview/:attachmentId` | Preview attachment file |

### Protected Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET/POST` | `/api/project/*` | Project CRUD operations |
| `GET/POST` | `/api/survey/*` | Survey settings and logic |
| `GET/POST` | `/api/answer/*` | Answer management |
| `GET/POST` | `/api/template/*` | Template CRUD with categories/tags |
| `GET/POST` | `/api/repo/*` | Question bank management |
| `GET/POST` | `/api/repo/userBook/*` | User's wrong/correct question book |
| `GET/POST` | `/api/repo/import` | Import questions from Excel |
| `GET` | `/api/repo/export` | Export questions to Excel |
| `GET/POST` | `/api/workflow/*` | Approval workflow operations |
| `GET` | `/api/ai/chat/models` | Get available AI models |
| `GET` | `/api/ai/chat/stream` | AI chat stream (SSE) |
| `GET/POST` | `/api/file/*` | File upload/download |
| `GET/POST` | `/api/dashboard/*` | Dashboard data |
| `GET/POST` | `/api/exercise/*` | Exercise history |
| `GET` | `/api/report/:shortId` | Report data |
| `GET/POST` | `/api/system/*` | System settings |

## Database

Tables use `t_` prefix. Core tables:

| Table | Description |
|-------|-------------|
| `t_project` | Surveys and exams |
| `t_answer` | Survey/exam responses |
| `t_template` | Question templates |
| `t_repo` | Question banks |
| `t_user_book` | User's question collection (wrong/correct) |
| `t_flow_operation` | Approval workflow history |
| `t_file` | Uploaded files |
| `t_user` | Users |
| `t_role` | Roles |
| `t_account` | Authentication accounts |
| `t_dashboard` | Dashboard configurations |
| `t_sys_info` | System settings |

## License

[MIT License](LICENSE)


