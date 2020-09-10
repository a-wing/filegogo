import Vue from 'vue'
import router from './router'

import { library } from '@fortawesome/fontawesome-svg-core';
import { faSearch, faArrowUp, faAngleRight, faAngleLeft, faAngleUp, faAngleDown, faUpload } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

library.add(faSearch, faArrowUp, faAngleRight, faAngleLeft, faAngleUp, faAngleDown, faUpload);
Vue.component('vue-fontawesome', FontAwesomeIcon);

import Buefy from 'buefy'
Vue.use(Buefy, {
  defaultIconComponent: 'vue-fontawesome',
  defaultIconPack: 'fas',
});

import 'buefy/dist/buefy.css'

import App from './App.vue'

new Vue({
  el: '#app',
  router,
  render: h => h(App)
})
