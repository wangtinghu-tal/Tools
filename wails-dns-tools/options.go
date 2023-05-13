package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

const FILE_NAME = "options.json"

type Options struct {
	TestHost string `json:"testHost"`

	TestDNSHost string `json:"testDNSHost"`
}

// 载入配置
func (opt *Options) load() {
	cwd, _ := os.Getwd()

	file, err := os.OpenFile(filepath.Join(cwd, FILE_NAME), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(content) == 0 {
		opt.TestHost = "www.baidu.com"
		data, _ := json.Marshal(opt)
		file.WriteString(string(data))
		return
	}

	json.Unmarshal(content, &opt)
}

// 保存配置文件
func (opt *Options) save() {
	cwd, _ := os.Getwd()
	file, err := os.OpenFile(filepath.Join(cwd, FILE_NAME), os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	data, _ := json.Marshal(opt)
	file.WriteString(string(data))
}

// 设置宽高值
func (opt *Options) set(key string, val string) {
	switch key {
	case "testHost":
		opt.TestHost = val
	case "testDNSHost":
		opt.TestDNSHost = val
	}
}

// SetData 存储新的配置
func (opt *Options) SetData(newOpt *Options) {
	if reflect.DeepEqual(newOpt, opt) {
		return
	}
	optionByte, _ := json.Marshal(newOpt)
	err := json.Unmarshal(optionByte, &opt)
	if err != nil {
		return
	}
	opt.save()
}

// GetData 获取配置信息
func (opt *Options) GetData() Options {
	return *opt
}
