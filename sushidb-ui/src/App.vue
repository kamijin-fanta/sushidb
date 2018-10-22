<template>
  <div id="app">
    {{store.counter}}
    <router-view class="view"></router-view>
    <button @click="incr()">incr</button>
  </div>
</template>

<script>
import { store } from './store'

export default {
  name: 'app',
  components: {
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
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
