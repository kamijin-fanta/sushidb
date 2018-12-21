// +build ignore

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/client-go/util/jsonpath"
	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"
	"time"
)

func main() {
	path := jsonpath.New("fuga")

	buf := new(bytes.Buffer)

	var data interface{}
	str := `[
		{"module": "hoge", "payload": {"lat": 123.456}},
		{"module": "fuga", "payload": {"lat": 123.456}}
	]`
	jsonErr := json.Unmarshal([]byte(str), &data)
	if jsonErr != nil {
		panic(jsonErr)
	}


	pathErr := path.Parse("{$[?(@.payload)]}")
	if pathErr != nil {
		panic(pathErr)
	}
	executeErr := path.Execute(buf, data)
	if executeErr != nil {
		panic(executeErr)
	}

	fmt.Printf("%+v\n", buf)

	return

	driver := tikv.Driver{}
	store, _ := driver.Open("tikv://192.168.75.79:2379")

	cli, err := tikv.NewRawKVClient([]string{"192.168.75.79:2379"}, config.Security{})

	key := []byte("transaction-test-key")

	tx, _ := store.Begin()
	tx.Set(key, []byte("test-data"))
	tx.Commit(context.Background())

	tx, _ = store.Begin()
	tx.Set(key, []byte("test-data-rollback"))
	tx.Rollback()

	tx, _ = store.Begin()
	res, _ := tx.Get(key)
	tx.Rollback()
	fmt.Printf("get data: %+v, \n", string(res))

	res2, _ := cli.Get(key)
	fmt.Printf("get data raw: %+v, \n", string(res2))

	cli.Put(key, []byte("from-raw-kv"))

	res3, _ := cli.Get(key)
	fmt.Printf("get data raw: %+v, \n", string(res3))

	time.Sleep(1000)

	tx, _ = store.Begin()
	res4, _ := tx.Get(key)
	fmt.Printf("get data: %+v, \n", string(res4))
	tx.Rollback()

	tx, _ = store.Begin()
	tx.Delete(key)
	tx.Commit(context.Background())

	tx, _ = store.Begin()
	res, _ = tx.Get(key)
	tx.Rollback()

	fmt.Printf("get data: %+v, \n", string(res))
	cli.Delete(key)

	return

	cli, err = tikv.NewRawKVClient([]string{"192.168.75.79:2379"}, config.Security{})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fmt.Printf("cluster ID: %d\n", cli.ClusterID())

	//key := []byte("tikv")
	//val := []byte("fooo")
	//
	//val, err = cli.Get(key)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("found val: %s for key: %s\n", val, key)

	keys, values, _ := cli.ReverseScan([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255}, 30)
	for i := range keys {
		fmt.Printf("found scan key: %s Value %s / HexKey: %x\n", keys[i], values[i], keys[i])
		//cli.Delete(keys[i])
	}
}
