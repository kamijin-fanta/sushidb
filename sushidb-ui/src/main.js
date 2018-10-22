import Vue from 'vue'
import App from './App.vue'
import VueRouter from 'vue-router';

Vue.config.productionTip = false

import Home from './components/Home.vue'
import HelloWorld from './components/HelloWorld.vue'

Vue.use(VueRouter)
const router = new VueRouter({
  mode: 'history',
  base: '/ui/',
  routes: [
    { path: '/', component: Home },
    { path: '/hello', component: HelloWorld },
  ]
})

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
