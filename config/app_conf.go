package config

import (
	"github.com/nie312122330/niexq-gotools/fileext"
	"gopkg.in/yaml.v2"
)

var AppConf *AppConfYaml

type AppConfYaml struct {
	// 数据库信息
	Server struct {
		AppName        string `yaml:"appName"`
		ServerPort     int    `yaml:"serverPort"`
		GinReleaseMode bool   `yaml:"ginReleaseMode"`
	} `yaml:"server"`
	//Chrome配置
	ChromeConf struct {
		ExecPath string `yaml:"execPath"`
		Pdfdir   string `yaml:"pdfdir"`
	} `yaml:"chromeConf"`
}

func init() {
	fileByte, err := fileext.ReadFileByte("app_conf.yaml")
	if nil != err {
		panic(err)
	}
	AppConf = &AppConfYaml{}
	err = yaml.Unmarshal(fileByte, AppConf)
	if nil != err {
		panic(err)
	}
}
