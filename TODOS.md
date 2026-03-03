# Frontend TODO List

This document tracks backend features that are missing frontend implementations.

## 1. Template Module (Complete Missing)

Create `website/src/api/template.ts`:

- [ ] `GET /template/list` - List templates with pagination
- [ ] `GET /template` - Get single template by ID
- [ ] `POST /template/create` - Create new template
- [ ] `POST /template/update` - Update existing template
- [ ] `POST /template/delete` - Delete template
- [ ] `GET /template/category/list` - List template categories
- [ ] `GET /template/tag/list` - List template tags

Create `website/src/types/template.ts`:

- [ ] Define `TemplateView`, `TemplateQuery`, `TemplateRequest` interfaces
- [ ] Define `CategoryView`, `TagView` interfaces

Views:

- [ ] Create `website/src/views/template/TemplateListView.vue` - Template list page
- [ ] Create `website/src/views/template/TemplateEditView.vue` - Template editor page

## 2. File Module (Complete Missing)

Create `website/src/api/file.ts`:

- [ ] `GET /file` - Get file by ID
- [ ] `GET /file/list` - List files (with projectId filter)
- [ ] `POST /file/upload` - Upload file
- [ ] `POST /file/delete` - Delete file
- [ ] `GET /file/downloadTemplate` - Download import template

Create `website/src/types/file.ts`:

- [ ] Define `FileView`, `FileQuery` interfaces

Views:

- [ ] Consider file manager component or integrate into existing views

## 3. Project Module Extensions

Add to `website/src/api/project.ts`:

- [ ] `GET /project/trash` - Get deleted projects (recycle bin)
- [ ] `POST /project/destroy` - Permanently delete project
- [ ] `GET /project/partner/list` - List project partners
- [ ] `POST /project/partner/create` - Add project partner
- [ ] `POST /project/partner/delete` - Remove project partner
- [ ] `GET /project/partner/download` - Export partners to Excel
- [ ] `POST /project/partner/import` - Import partners from Excel
- [ ] `POST /project/selectUser` - Select user for editor picker
- [ ] `POST /project/selectRole` - Select role for editor picker
- [ ] `POST /project/selectTemplate` - Select template for editor picker

Add to `website/src/types/project.ts`:

- [ ] Define `ProjectPartnerView`, `ProjectPartnerQuery`, `ProjectPartnerRequest`
- [ ] Define `SelectUserView`, `SelectRoleView`, `SelectTemplateView`

Views:

- [ ] Add recycle bin tab to ProjectListView
- [ ] Create partner management component/dialog in ProjectEditView
- [ ] Implement partner import/export functionality

## 4. Answer Module Extensions

Add to `website/src/api/answer.ts`:

- [ ] `GET /answer/trash` - Get deleted answers (recycle bin)
- [ ] `POST /answer/create` - Manually create answer
- [ ] `POST /answer/update` - Manually update answer
- [ ] `POST /answer/destroy` - Permanently delete answer
- [ ] `POST /answer/upload` - Import answers from Excel

Add to `website/src/types/answer.ts`:

- [ ] Define `AnswerUploadResult` interface

Views:

- [ ] Add recycle bin tab to AnswerListView
- [ ] Add answer import functionality
- [ ] Add manual answer creation/editing (admin feature)

## 5. Survey Module Extensions

Add to `website/src/api/survey.ts`:

- [ ] `POST /public/validateProject` - Validate survey access
- [ ] `POST /public/statistics` - Get survey statistics
- [ ] `GET /public/preview/:attachmentId` - Preview attachment
- [ ] `POST /public/loadQuery` - Load query data
- [ ] `POST /public/getQueryResult` - Get query result
- [ ] `POST /public/loadDict` - Load dictionary data
- [ ] `POST /public/loadExamResult` - Load exam result
- [ ] `POST /public/loadLinkResult` - Load link result
- [ ] `GET /api/survey/setting` - Get survey settings
- [ ] `POST /api/survey/setting` - Update survey settings
- [ ] `GET /api/survey/logic` - Get survey logic rules
- [ ] `POST /api/survey/logic` - Update survey logic rules

Add to `website/src/types/survey.ts`:

- [ ] Define `SurveyStatistics`, `SurveySetting`, `SurveyLogic` interfaces
- [ ] Define `QueryRequest`, `QueryResult`, `DictRequest`, `DictResult`
- [ ] Define `ExamResult`, `LinkResult` interfaces

Views:

- [ ] Add statistics view component
- [ ] Add survey settings management UI
- [ ] Add survey logic editor component
- [ ] Support query/dictionary question types

## 6. User Module Extensions

Add to `website/src/api/auth.ts`:

- [ ] `GET /listUserTask` - Get current user's tasks
- [ ] `GET /listHistoryTask` - Get current user's history tasks

Add to `website/src/types/user.ts`:

- [ ] Define `UserTaskView`, `UserTaskQuery` interfaces

Views:

- [ ] Create user task center view (optional, depends on UX requirements)

## 7. System Module Extensions

Add to `website/src/api/system.ts`:

- [ ] `POST /system/update` - Update system info
- [ ] `GET /system/aiSetting` - Get AI settings
- [ ] `GET /system/user/list` - List users (system admin)
- [ ] `POST /system/user/create` - Create user (system admin)
- [ ] `POST /system/user/update` - Update user (system admin)
- [ ] `POST /system/user/delete` - Delete user (system admin)

Add to `website/src/types/user.ts` or `website/src/types/system.ts`:

- [ ] Define `SystemInfoRequest`, `AISetting` interfaces

Views:

- [ ] Add system settings page for admin
- [ ] Integrate AI settings configuration
- [ ] Ensure UserListView uses appropriate API endpoints

## 8. Router Updates

Update `website/src/router/index.ts`:

- [ ] Add template management routes
- [ ] Add file manager routes (if creating standalone view)
- [ ] Add user task routes (if implementing task center)
- [ ] Add recycle bin routes

## 9. Navigation Updates

Update `website/src/components/common/Layout.vue`:

- [ ] Add template management menu item
- [ ] Add file manager menu item (if applicable)
- [ ] Add recycle bin access

## Priority Order

1. **High Priority** - Core functionality gaps
   - Template module (needed for survey creation flow)
   - Project partner management (essential for survey distribution)
   - Recycle bin for projects and answers

2. **Medium Priority** - Enhanced features
   - File management API
   - Survey statistics and settings
   - Answer import functionality

3. **Low Priority** - Nice to have
   - User task center
   - AI settings
   - Query/dictionary question types
   - Exam/link result features
