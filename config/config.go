package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// MirrorConfig 静态服务器参数
type MirrorConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	MirrorURI string `yaml:"mirror_uri"`
}

// NewMirrorConfig 初始化默认参数
func NewMirrorConfig() MirrorConfig {
	return MirrorConfig{
		Port:      9999,
		MirrorURI: "/mirror",
		Host:      "127.0.0.1",
	}
}

// Load 从配置文件加载参数
func (m *MirrorConfig) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("加载文件配置文件失败，原因%s\n", err.Error())
		return err
	}

	var mc MirrorConfig
	err = yaml.Unmarshal(data, &mc)
	if err != nil {
		log.Printf("解析配置文件失败，原因%s\n", err.Error())
		return err
	}

	// 加载port
	if mc.Port != 0 {
		m.Port = mc.Port
	}

	// 加载URI
	if mc.MirrorURI != "" {
		m.MirrorURI = mc.MirrorURI
	}

	if mc.Host != "" {
		m.Host = mc.Host
	}
	return nil
}
