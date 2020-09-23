import Vue from 'vue'

import element from '@/utils/element'
import 'element-ui/lib/theme-chalk/index.css'

import '@/styles/index.scss' // global css

import App from '@/App'
import store from '@/store/index'

import 'xterm/css/xterm.css'

Vue.use(element)

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
    el: '#app',
    store,
    render: h => h(App)
})
