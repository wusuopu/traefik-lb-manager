import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import RuleView from '../views/RuleView.vue'
import CertificateView from '../views/CertificateView.vue'

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
      component: CertificateView,
    },
  ],
})

export default router
