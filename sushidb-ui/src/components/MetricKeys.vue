<template>
  <div class="metric-keys">
    <h1>metrics</h1>
    <at-table :columns="tableColumns" :data="tableData"></at-table>
  </div>
</template>

<script>
import {
  Table as AtTable,
  Button as AtButton,
} from 'at-ui'

export default {
  name: 'metric-keys',
  components: {
    AtTable,
  },
  mounted () {
    this.fetch();
  },
  computed: {
    tableData () {
      return this.keys.map((value, index) => ({ id: index + 1, key: value }))
    }
  },
  data () {
    return {
      keys: [],
      tableColumns: [
        { title: '#', key: 'id'},
        { title: 'Metric Key', key: 'key'},
        {
          title: 'Operation',
          render: (h, params) => {
            return h('div', [
              h(AtButton, {
                props: {
                  size: 'small',
                  hollow: true
                },
                style: {
                  marginRight: '8px'
                },
                on: {
                  click: () => {
                    console.log('view', params)
                  }
                }
              }, 'View Metrics'),
            ])
          }
        }
      ]
    }
  },
  methods: {
    fetch() {
      fetch('/metric/keys')
        .then(res => res.json())
        .then(res => { this.keys = res })
    }
  }
}
</script>

<style scoped>
</style>
