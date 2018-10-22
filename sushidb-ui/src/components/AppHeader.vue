<template>
  <div class="app-header">
    <h1>SushiDB</h1>
    <at-menu mode="horizontal" active-name="1" router>
      <at-menu-item :to="{ name: 'home' }"><i class="icon icon-home"></i>Home</at-menu-item>
      <at-menu-item :to="{ name: 'metric-keys' }"><i class="icon icon-layers"></i>Kyes</at-menu-item>
    </at-menu>
  </div>
</template>

<script>
import {
  Menu as AtMenu,
  MenuItem as AtMenuItem,
  MenuItemGroup as AtMenuItemGroup,
  Submenu as AtSubmenu,
} from 'at-ui'

export default {
  name: 'app-header',
  components: {
    AtMenu,
    AtMenuItem,
    AtMenuItemGroup,
    AtSubmenu,
  },
  data () {
    return {
    }
  },
  methods: {
  }
}
</script>

<style scoped>
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2rem;
  border-bottom: 1px solid #e2ecf4;
}
h1 {
  min-width: min-content;
    color: #6582c7;
    font-weight: normal;
}
.at-menu {
  border: none;
}
</style>
