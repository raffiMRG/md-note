import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue'),
    meta: { guestOnly: true, hideNavBar: true },
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('../views/RegisterView.vue'),
    meta: { guestOnly: true, hideNavBar: true },
  },
  {
    path: '/',
    name: 'dashboard',
    component: () => import('../views/DashboardView.vue'),
  },
  {
    path: '/notes/new',
    name: 'note-new',
    component: () => import('../views/NoteEditorView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/notes/:id/edit',
    name: 'note-edit',
    component: () => import('../views/NoteEditorView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/notes/:id',
    name: 'note-detail',
    component: () => import('../views/NoteDetailView.vue'),
    meta: { notepadLayout: true },
  },
  {
    path: '/tags',
    name: 'tags',
    component: () => import('../views/TagManagerView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('../views/SettingsView.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
  },
  {
    path: '/admin/users',
    name: 'admin-users',
    component: () => import('../views/AdminUsersView.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    return { name: 'dashboard' }
  }
  if (to.meta.guestOnly && authStore.isAuthenticated) {
    return { name: 'dashboard' }
  }
})

export default router
