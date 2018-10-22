<template>
  <div id="app">
    <app-header></app-header>
    <router-view class="view"></router-view>
    <!-- {{store.counter}}
    <button @click="incr()">incr</button> -->
  </div>
</template>

<script>
import { store } from './store'
import AppHeader from './components/AppHeader'

export default {
  name: 'app',
  components: {
    AppHeader,
  },
  data () {
    return {
      store,
    }
  },
  methods: {
    incr() {
      this.store.counter++;
    }
  }
}
</script>

<style>
#app {
  font-family: 'Yu Gothic Medium', YuGothic, YuGothic M, meiryo, Helvetica Neue, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
#app .view {
  margin: 2rem;
}

#app .view h1 {
  margin: 1rem 0;
}
</style>
