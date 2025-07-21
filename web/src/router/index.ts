import Home from '@/views/Home.vue';
import Library from '@/views/Library.vue';
import Main from '@/views/Main.vue';
import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/home',
    },
    {
      path: '/',
      component: Main,
      children: [
        {
          path: '/home',
          name: 'home',
          component: Home,
        },
        {
          path: '/library/:id?/:page?',
          name: 'library',
          component: Library,
        },
        {
          path: '/series/:id/:page?',
          name: 'series',
          component: () => import('@/views/Series.vue'),
        },
        {
          path: '/tasks',
          name: 'tasks',
          component: () => import('@/views/Tasks.vue'),
        },
      ],
    },
    {
      path: '/series/:seriesID/comic/:id/:page',
      name: 'comic',
      component: () => import('@/views/Comic.vue'),
    },
    {
      path: '/series/:seriesID/novel/:id/:page',
      name: 'novel',
      component: () => import('@/views/Novel.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/Register.vue'),
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('@/views/Settings/Settings.vue'),
      redirect: '/settings/about',
      children: [
        {
          path: '/settings/user',
          name: 'settings-user',
          component: () => import('@/views/Settings/User.vue'),
        },
        {
          path: '/settings/library',
          name: 'settings-library',
          component: () => import('@/views/Settings/Library.vue'),
        },
        {
          path: '/settings/user-manage',
          name: 'settings-user-manage',
          component: () => import('@/views/Settings/UserManage.vue'),
        },
        {
          path: '/settings/about',
          name: 'settings-about',
          component: () => import('@/views/Settings/About.vue'),
        },
      ],
    },
  ],
});

// 添加全局前置守卫，未认证跳转 /login
router.beforeEach((to, from, next) => {
  const publicPages = ['/login', '/register'];
  const authRequired = !publicPages.includes(to.path);
  const token = localStorage.getItem('token');
  if (authRequired && !token) {
    next('/login');
  } else {
    next();
  }
});

export default router;
