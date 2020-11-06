package main

import (
	"fmt"
	"github.com/simonhylander/booker"
)

func main() {

	/*cfgPath := flag.String("p", "./cmd/api/conf.local.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)
	checkErr(api.Start(cfg))*/

	uri := "bolt://localhost:7687"
	username := "neo4j"
	password := "booking"
	str, err := booker.HelloWorld(uri, username, password, false)
	fmt.Println(str)
	fmt.Println(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}