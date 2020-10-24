import Vue from 'vue'
import Router from 'vue-router'

import Overview from './pages/overview.vue'
Vue.use(Router)

const routes = [
  {
    path: '/',
    component: Overview
  },
  {
    path: '/t/:id',
    component: Overview
  }
]

export default new Router({
  routes
})
