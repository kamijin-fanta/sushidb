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
      return this.keys.map((value, index) => ({
        id: index + 1,
        key: value.metric_id,
        type: value.type,
      }))
    }
  },
  data () {
    return {
      keys: [],
      tableColumns: [
        { title: '#', key: 'id'},
        { title: 'Metric Key', key: 'key' },
        { title: 'Type', key: 'type' },
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
                    switch (params.item.type) {
                      case 'message':
                        this.$router.push({ name: 'message-query', params: { metric_id: params.item.key }});
                        break;
                      case 'single':
                        this.$router.push({ name: 'single-query', params: { metric_id: params.item.key }});
                        break;
                    }
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
      fetch(`/keys`)
        .then(res => res.json())
        .then(res => { this.keys = res })
    }
  }
}
</script>

<style scoped>
.selection {
  margin: 2rem 0;
}
</style>
