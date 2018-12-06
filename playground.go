// +build ignore

package main

import (
	"fmt"
	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"
)

func main() {
	cli, err := tikv.NewRawKVClient([]string{"192.168.75.79:2379"}, config.Security{})
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
