<template>
    <div class="query-single">
      <h1>Query Single Data</h1>
      <div class="forms">
        <div>
          <at-input v-model="inputMetricId"></at-input>
        </div>
        <div>
          <at-radio v-model="sort" label="asc">ASC</at-radio>
          <at-radio v-model="sort" label="desc">DESC</at-radio>
        </div>
      </div>
      <at-table :columns="tableColumns" :data="tableData"></at-table>
    </div>
  </template>
  
  <script>
  import {
    Table as AtTable,
    Input as AtInput,
    Radio as AtRadio,
  } from 'at-ui'
  
  export default {
    name: 'query-single',
    components: {
      AtTable,
      AtInput,
      AtRadio,
    },
    mounted () {
      this.fetch();
      this.inputMetricId = this.metricId;
    },
    computed: {
      tableData () {
        return this.response.rows && this.response.rows.map((row, index) => ({
          id: index + 1,
          time: row.time,
          value: row.value,
        }))
      },
      metricId () {
        return this.$route.params.metric_id
      },
    },
    watch: {
      metricId () {
        this.fetch();
      },
      sort () {
        this.fetch();
      },
    },
    data () {
      return {
        response: {},
        tableColumns: [
          { title: '#', key: 'id'},
          { title: 'Time', key: 'time' },
          { title: 'Value', key: 'value' },
        ],

        // inputs
        inputMetricId: '',
        sort: 'desc',
      }
    },
    methods: {
      fetch() {
        fetch(`/metric/single/${this.metricId}?sort=${this.sort}`)
          .then(res => res.json())
          .then(res => { this.response = res })
      }
    }
  }
  </script>
  
  <style scoped>
  .forms {
    margin: 2rem auto;
    max-width: 600px;
    
  }
  </style>
  