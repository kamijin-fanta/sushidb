<template>
  <div class="home">
    home panel ({{store.counter}})
  </div>
</template>

<script>
import {store} from '../store'

export default {
  name: 'home',
  props: {
  },
  data () {
    return {
      store,
    }
  },
}
</script>
<style scoped>
</style>
