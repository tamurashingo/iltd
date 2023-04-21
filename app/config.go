// The MIT License (MIT)
//
// Copyright (c) 2015 tamura shingo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Config struct {
	Appname *string       `json:"appname"`
	Server  *ServerConfig `json:"server"`
	Db      *DbConfig     `json:"db"`
}

type ServerConfig struct {
	Port *string `json:"port"`
}

type DbConfig struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Host     *string `json:"host"`
	Port     *string `json:"port"`
	Dbname   *string `json:"dbname"`
}

func configLoad(filename string) (*Config, error) {
	var config Config
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(file, &config)

	return &config, nil
}

func configValidate(config *Config) error {
	if *config.Appname == "" {
		return errors.New("AppNnameが未設定です")
	}

	if *config.Server.Port == "" {
		return errors.New("Server.Portが未設定です")
	}

	if *config.Db.Username == "" {
		return errors.New("Db.Usernameが未設定です")
	}
	if *config.Db.Password == "" {
		return errors.New("Db.Passwordが未設定です")
	}
	if *config.Db.Host == "" {
		return errors.New("Db.Hostが未設定です")
	}
	if *config.Db.Port == "" {
		return errors.New("Db.Portが未設定です")
	}
	if *config.Db.Dbname == "" {
		return errors.New("Db.Dbnameが未設定です")
	}

	return nil
}
