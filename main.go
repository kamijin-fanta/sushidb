package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	cli, err := tikv.NewRawKVClient(strings.Split(pdAddress, ","), config.Security{})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fmt.Printf("cluster ID: %d\n", cli.ClusterID())

	r := gin.Default()
	ApiServer(r, cli)
	UiServer(r)
	r.Run()
}
