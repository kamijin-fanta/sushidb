# sushidb

## development

```
go mod download
go get github.com/codegangsta/gin
gin --port 3006
```

## API

### GET /ping

### GET /cluster

### POST /metric/{single|message}/:id/:time

- id: key name (example: hoge)
- time: nano second time (example: 1544068003882000)
  - fail case: time < 1000000000000000 || time > 9000000000000000

```bash
$ curl -XPOST localhost:3000/metric/single/hoge/1544068003882000 -d '{"app": "hoge", "la": 0.24}'
{"ok":1}
```

### GET /metric/{single|message}/:id?lower={ns_time}&upper={ns_time}&limit={num}&sort={asc|desc}

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

### POST /metric/{single|message}/:id

advanced quering api

sample request to `POST /metric/single/piyo`

```json
{
  "filters": [
    {
      "type": "lte",
      "path": "$",
      "value": 3
    }
  ],
  "limit": 1000,
  "max_skip": 1000,
  "cursor": 1544068003885000
}
```

response

```json
{
  "metric_id": "piyo",
  "rows": [
    {
      "time": 1544068003882000,
      "value": 2.22
    },
    {
      "time": 1544068003881000,
      "value": 1.11
    }
  ],
  "query_time_ns": 33867300,
  "cursor": 1544068003881000
}
```


### GET /keys/

## UI

```bash
$ chrome localhost:3000/ui 
```

## キー設計

- フォーマット
  - `[prefix 2bytes]_[metricKey some bytes]_[subtype 1 byte]_[time ns 8 bytes]`
  - 各項目はアンダーバー(0x5f)で区切る
- prefix: 値の種別・バージョンが入る
- metricKey: キー名などが入る
- subtype: 圧縮後の解像度など、該当のキーへの補助的な種別が入る
- time: ビッグエンディアンのint64値として、ナノ秒を格納する


### Prefix

#### s1

- 値を格納する
- subtype: Resolution
- body: msgpackでマーシャルされた単一の値

#### m1

- メッセージを格納する
- subtype: Resolution
- body: msgpackでマーシャルされた値

#### k1

- キーのリストを格納する
- subtype: prefix type
  - 0: v1
  - 1: m1
- body: empty


### Subtype

#### Resolution

- 0: 生
- 1: 解像度を維持しつつデータ数を抑える
- 2: 1分毎に丸める
- 3: 1時間毎に丸める
- 4: 1日毎に丸める
