import Vue from 'vue'
import router from './router'

import { library } from '@fortawesome/fontawesome-svg-core'
import { faUpload, faDownload } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

import Buefy from 'buefy'
import 'buefy/dist/buefy.css'

import App from './App.vue'

library.add(faUpload, faDownload)
Vue.component('vue-fontawesome', FontAwesomeIcon)

Vue.use(Buefy, {
  defaultIconComponent: 'vue-fontawesome',
  defaultIconPack: 'fas'
})

const el = document.createElement('div')
document.body.appendChild(el)
new Vue({ // eslint-disable-line no-new
  el: el,
  router,
  extends: App
})
