import { createRouter, createWebHashHistory } from 'vue-router'
import ChatLayout from '@/views/chat/index.vue'

const routes = [
  {
    path: '/',
    name: 'Root',
    component: ChatLayout,
    redirect: '/chat',
    children: [
      {
        path: '/chat',
        name: 'Chat',
        component: () => import('../views/chat/index.vue'),
      },
    ],
  },
  // {
  //   path: '/404',
  //   name: '404',
  //   component: () => import('../views/exception/404.vue'),
  // },

  // {
  //   path: '/500',
  //   name: '500',
  //   component: () => import('@/views/exception/500.vue'),
  // },
  {
    path: '/:pathMatch(.*)*',
    name: 'notFound',
    redirect: '/404',
  },
]

export const router = createRouter({
  history: createWebHashHistory('/kael/'),
  routes,
  scrollBehavior: () => ({ left: 0, top: 0 }),
})

export async function setupRouter(app) {
  app.use(router)
  await router.isReady()
}
