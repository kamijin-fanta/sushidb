# sushidb

## development

```
go mod download
go get github.com/codegangsta/gin
gin
```

## API

### GET /ping

### GET /cluster

### POST /metric/single/:id/:time

- id: key name (example: hoge)
- time: nano second time (example: 1544068003882000)
  - fail case: time < 1000000000000000 || time > 9000000000000000

```bash
$ curl -XPOST localhost:3000/metric/single/hoge/1544068003882000 -d '{"app": "hoge", "la": 0.24}'
{"ok":1}
```

### GET /metric/single/:id?lower={ns_time}&upper={ns_time}&limit={num}&sort={asc|desc}

- id: key name
  - format: string
- lower: lower limit of fetch time range
  - default: none spec
  - format: nano second time
- upper: upper limit of fetch time range
  - default: none spec
  - format: nano second time
- sort: Direction of fetch
  - default: desc (the latest data is the first)
  - format: string. asc or desc

```json
$ curl localhost:3000/metric/single/hoge
{
  "metric_id":"hoge",
  "rows":[
    {
      "time":1544068003882000,
      "value":{"app":"hoge","la":0.24}
    },
    {
      "time":1544068003884000,
      "value":{"app":"hoge","la":0.26}
    },
    {
      "time":1544068003883000,
      "value":{"app":"hoge","la":0.24}
    }
  ]
}
```

### GET /metric/keys

## UI

```bash
$ chrome localhost:3000/ui 
```
