package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
)

const FILE_NAME = "options.json"

type Options struct {
	PDFHeight     string `json:"pdfHeight"`
	PDFWidth      string `json:"pdfWidth"`
	PDFExportPath string `json:"pdfExportPath"`
}

type React struct {
	Regexp string `json:"regexp"`
	Value  string `json:"value"`
}

type EachList struct {
	regexp *regexp.Regexp
	value  string
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
		opt.PDFHeight = "768"
		opt.PDFWidth = "1024"
		opt.PDFExportPath = "./"
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
	case "pdfHeight":
		opt.PDFHeight = val
	case "pdfWidth":
		opt.PDFWidth = val

	case "pdfExportPath":
		opt.PDFExportPath = val
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
