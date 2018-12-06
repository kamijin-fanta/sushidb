import Vue from 'vue'
import App from './App.vue'
import VueRouter from 'vue-router';

import 'at-ui-style'
// import AtComponents from 'at-ui'
// Vue.use(AtComponents)

Vue.config.productionTip = false

import Home from './components/Home.vue'
import HelloWorld from './components/HelloWorld.vue'
import MetricKeys from './components/MetricKeys.vue'
import QuerySingle from './components/QuerySingle.vue'

Vue.use(VueRouter)
const router = new VueRouter({
  mode: 'history',
  base: '/ui/',
  routes: [
    { path: '/', redirect: { name: 'home' } },
    { path: '/home', name: 'home', component: Home },
    { path: '/hello', name: 'hello', component: HelloWorld },
    { path: '/metric/keys', name: 'metric-keys', component: MetricKeys },
    { path: '/metric/single/:metric_id', name: 'single-query', component: QuerySingle },
  ]
})

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
