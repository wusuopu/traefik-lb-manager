import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import RuleView from '../views/RuleView.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: HomeView,
    },
    {
      path: '/workspaces/:id/rules',
      name: 'Rules',
      component: RuleView,
    },
    {
      path: '/workspaces/:id/authentications',
      name: 'Authentications',
      component: RuleView,
    },
    {
      path: '/workspaces/:id/certificates',
      name: 'Certificates',
      component: RuleView,
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
  ],
})

export default router
