package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type user struct {
	Email string
}
type appConf struct {
	NsMaxAge  float64         `yaml:nsmaxage`
	Users     map[string]user `yaml:users`
	CheckFreq int64           `yaml:checkfreq`
}

func NewConf() appConf {
	data, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		panic(err.Error())
	}
	conf := appConf{}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		panic(err.Error())
	}
	return conf
}
