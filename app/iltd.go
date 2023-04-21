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
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
)

type IltdApp struct {
	config *Config
}

// 指定した設定ファイルを元に新規インスタンスを生成する
func IltdNew(configfile string) (*IltdApp, error) {
	config, err := configLoad(configfile)
	if err != nil {
		return nil, err
	}

	err = configValidate(config)
	if err != nil {
		return nil, err
	}

	return &IltdApp{
		config: config,
	}, nil
}

// アプリケーションを起動する
// * ルーティングの設定
// * サーバ起動
func (iltdApp *IltdApp) IltdRun() error {
	log.Debug("Run")
	_, err := dbInit(iltdApp.config)
	if err != nil {
		return err
	}
	defer dbClose()

	router := mux.NewRouter().StrictSlash(true)

	// アプリ名称取得
	appRouter(router, iltdApp)

	// Task検索
	searchRouter(router, iltdApp)
	// 登録
	registRouter(router, iltdApp)
	// 更新
	updateRouter(router, iltdApp)

	// ユーザ一覧取得
	userRouter(router, iltdApp)

	// ファイル
	staticRouter(router, iltdApp)

	log.Info("listening :", *iltdApp.config.Server.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", *iltdApp.config.Server.Port), router)
	return err
}

func appRouter(r *mux.Router, iltdApp *IltdApp) *mux.Router {
	r.
		Methods("GET").
		Path("/api/app/name").
		Name("App API").
		Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info(r.URL.Path)
			result := fmt.Sprintf("{\"appname\": \"%s\"}", *iltdApp.config.Appname)
			successResponse([]byte(result), w)
		}))

	return r
}

// 検索API
//   /api/search
//   /api/search?user=xxxx
//   /api/search?update=yyyyMMdd
//   /api/search?user=xxxx&update=yyyyMMdd
//
func searchRouter(r *mux.Router, iltdApp *IltdApp) *mux.Router {
	r.
		Methods("GET").
		Path("/api/search").
		Name("Search API").
		Handler(http.HandlerFunc(showtaskController))

	return r
}

// 登録API
//   /api/task
//   {
//     "type": "",
//     "name": "",
//     "plan": "",
//     "result": "",
//     "unit": "",
//     "due": "",
//     "person": ""
//   }
func registRouter(r *mux.Router, iltdApp *IltdApp) *mux.Router {
	r.
		Methods("POST").
		Path("/api/task").
		Name("Regist API").
		Handler(http.HandlerFunc(registtaskController))

	return r
}

// 更新API
// * ステータス
//   /api/task/id
//   {
//     "status": ""
//   }
//
// * 詳細
//   /api/task/id
//   {
//     "type": "",
//     "name": "",
//     "plan": "",
//     "result": "",
//     "unit": "",
//     "due": "",
//     "person": ""
//   }
func updateRouter(r *mux.Router, iltdApp *IltdApp) *mux.Router {
	r.
		Methods("PUT").
		Path("/api/task/{id}").
		Name("Update API").
		Handler(http.HandlerFunc(updatetaskController))

	return r
}

func userRouter(r *mux.Router, iltdApp *IltdApp) *mux.Router {
	r.
		Methods("GET").
		Path("/api/user").
		Name("User API").
		Handler(http.HandlerFunc(alluserController))

	return r
}

// HTML等のファイル
func staticRouter(r *mux.Router, iltdApp *IltdApp) *mux.Router {
	r.
		Methods("GET").
		PathPrefix("/").
		Name("Static Files").
		Handler(http.FileServer(http.Dir("./static")))

	return r
}
