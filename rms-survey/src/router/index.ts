import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { public: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { public: true }
  },
  {
    path: '/s/:code',
    name: 'SurveyFill',
    component: () => import('@/views/survey/SurveyFillView.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    component: () => import('@/components/common/Layout.vue'),
    children: [
      {
        path: '',
        redirect: '/project'
      },
      {
        path: 'project',
        name: 'ProjectList',
        component: () => import('@/views/project/ProjectListView.vue')
      },
      {
        path: 'project/create',
        name: 'ProjectCreate',
        component: () => import('@/views/project/ProjectEditView.vue')
      },
      {
        path: 'project/:id/edit',
        name: 'ProjectEdit',
        component: () => import('@/views/project/ProjectEditView.vue')
      },
      {
        path: 'answer',
        name: 'AnswerList',
        component: () => import('@/views/answer/AnswerListView.vue')
      },
      {
        path: 'answer/:id',
        name: 'AnswerDetail',
        component: () => import('@/views/answer/AnswerDetailView.vue')
      },
      {
        path: 'system/user',
        name: 'UserList',
        component: () => import('@/views/system/UserListView.vue')
      },
      {
        path: 'system/role',
        name: 'RoleList',
        component: () => import('@/views/system/RoleListView.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  // Public routes
  if (to.meta.public) {
    next()
    return
  }

  // Check authentication
  if (!authStore.token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  // Fetch user if not loaded
  if (!authStore.user) {
    try {
      await authStore.fetchCurrentUser()
    } catch {
      next({ name: 'Login', query: { redirect: to.fullPath } })
      return
    }
  }

  next()
})

export default router
