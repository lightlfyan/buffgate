package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var Config = &CfgType{}

func loadCfg(path string, cfg interface{}) {
	bins, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(bins, cfg)
	if err != nil {
		fmt.Println("err: ", err)
		panic(err)
	}
}

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("err: ", err)
		panic(err)
	}
	loadCfg(dir+"/config/config.json", Config)
}
