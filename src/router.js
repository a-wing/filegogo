import Vue from 'vue';
import Router from 'vue-router';
Vue.use(Router);

import Overview from './pages/overview.vue';

const routes = [
  {
    path: '/',
    component: Overview
  },
  {
    path: '/t/:id',
    component: Overview
  },
]

export default new Router({
  routes
});
