# TODOs

## Before First Run

- [ ] Set database DSN in `configs/config.yaml` (currently placeholder)
- [ ] Set `JWT_SECRET` env var (min 64 chars for HS512)
- [ ] Set `AI_API_KEY` env var for SiliconFlow integration
- [ ] Create upload directory: `mkdir -p uploads`

## Incomplete Business Logic

### Workflow (internal/service/flow_service.go)
- [ ] `ApprovalTask`: state transition is not persisted — needs to update answer status in DB atomically
- [ ] `GetAuditRecord`: returns empty list — needs a flow_operation table or equivalent to record approval history
- [ ] `GetRevertNodes`: returns empty list — needs to query previous approval nodes
- [ ] `SaveFlow` / `Deploy`: flow schema save and deployment logic not implemented
- [ ] `GetTaskInfo`: endpoint exists in Java (`/getTaskInfo`) but not wired in Go router
- [ ] `LoadSchema`: endpoint exists in Java (`/loadSchema`) but not wired in Go router
- [ ] Flow state machine: pending → running → approved/rejected/cancelled transitions need DB transaction wrapping

### Report (internal/service/report_service.go)
- [ ] `GetData`: only returns total answer count — Java version returns per-question statistics (choice distributions, text answers, etc.)
- [ ] Report data needs to parse survey JSON structure and aggregate answers per question

### Dashboard (internal/service/dashboard_service.go)
- [ ] `SaveDashboard` / `UpdateDashboard`: only List is implemented, no create/update/delete
- [ ] Dashboard handler only has GET — Java DashboardApi likely has PUT for saving dashboard config

### Exercise (internal/service/exercise_service.go)
- [ ] Exercise scoring logic not implemented (exam score calculation)
- [ ] `GetExerciseDetail`: only list is implemented, no detail view

## Router Gaps

- [ ] Flow handler routes (flow.go, ai.go, file.go, repo.go, template.go, exercise.go, dashboard.go, report.go) are manually wired in router.go but need verification against Java endpoint paths
- [ ] `/api/workflow/loadSchema` and `/api/workflow/getTaskInfo` not registered
- [ ] Verify all 126 Java endpoints are covered (current Go router covers ~90)

## Known Simplifications

- [ ] `flow_service.GetFlowEntry`: returns hardcoded `status: "active"` — needs real status from DB
- [ ] `flow_service.GetFlowTasks`: maps `ExamExerciseType` as status — incorrect, needs dedicated flow status field
- [ ] `flow_service.Statics`: counts by `ExamExerciseType` — should count by actual flow approval status
- [ ] Answer export: Excel column headers are hardcoded — should be dynamic based on survey question structure
- [ ] Data permission filtering: only checks `create_by` and project_partner — verify matches Java DataPermAspect logic exactly

## Not Implemented

- [ ] `UserBook` (address book) endpoints — DTOs exist but handler routes may be missing
- [ ] `RepoTemplate` import/export — skeleton only
- [ ] Captcha: uses simple in-memory map, no TTL cleanup goroutine (memory leak on high traffic)
- [ ] RSA key rotation / management endpoint
