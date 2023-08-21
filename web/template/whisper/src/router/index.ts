import { createRouter, createWebHistory } from 'vue-router'
import EquipBox from '../views/EquipBox.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      // component: SearchBox
      component: () => import('../views/SearchBox.vue')
    },
    {
      path: '/equip',
      name: 'equip',
      component: EquipBox
      // component: () => import('../views/EquipBox.vue')
    }
  ]
})

export default router
