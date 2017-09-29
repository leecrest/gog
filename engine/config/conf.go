package config

// ini配置文件加载

import (
	"bufio"
	"io"
	"os"
	"strings"
	"strconv"
)

const middle = "$"

type INIConfig struct {
	data map[string]string
	section string
}

func INILoad(path string) (*INIConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg INIConfig
	cfg.data = make(map[string]string)
	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			cfg.section = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(cfg.section) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		first := strings.TrimSpace(s[:index])
		if len(first) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := cfg.section + middle + first
		cfg.data[key] = strings.TrimSpace(second)
	}
	return &cfg, nil
}

func (cfg *INIConfig) Read(section string, key string, def string) (string) {
	var k = section + middle + key
	var v, found = cfg.data[k]
	if !found {
		return def
	}
	return v
}

func (cfg *INIConfig) ReadInt(section string, key string, def int) (int) {
	var k = section + middle + key
	var v, found = cfg.data[k]
	if !found {
		return def
	}
	var val, err = strconv.Atoi(v)
	if err != nil {
		panic(err)
		return def
	}
	return val
}

func (cfg *INIConfig) ReadUint32(section string, key string, def uint32) (uint32) {
	var k = section + middle + key
	var v, found = cfg.data[k]
	if !found {
		return def
	}
	var val, err = strconv.Atoi(v)
	if err != nil {
		panic(err)
		return def
	}
	return uint32(val)
}
