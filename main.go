package main

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pingcap/pd/client"
	"log"
	"os"
	"strings"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}
	pdAddress := os.Getenv("PD_ADDRESS")

	addressList := strings.Split(pdAddress, ",")
	rawClient, err := tikv.NewRawKVClient(addressList, config.Security{})
	if err != nil {
		panic(err)
	}
	defer rawClient.Close()

	pdClient, err := pd.NewClient(addressList, pd.SecurityOption{})
	if err != nil {
		panic(err)
	}
	defer pdClient.Close()

	store := Store{
		rawKvClient: *rawClient,
		pbClient:    pdClient,
	}

	fmt.Printf("cluster ID: %d\n", rawClient.ClusterID())

	r := gin.Default()
	pprof.Register(r) // enabled /debug/pprof/

	ApiServer(r, &store)
	UiServer(r)

	r.Run()
}
