import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '../views/HomePage.vue'
import SearchBox from '../views/SearchBox.vue'
import EquipBox from '../views/EquipBox.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomePage
    },
    {
      path: '/search',
      name: 'search',
      component: SearchBox
    },
    {
      path: '/equip',
      name: 'equip',
      component: EquipBox
    }
  ]
})

export default router
