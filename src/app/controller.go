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
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

const (
	error_response = `
{"result": "false",
"message": "%s"}
`
)

func showtaskController(w http.ResponseWriter, r *http.Request) {
	log.Info(r.URL.Path)
	request := map[string]string{
		"user":   r.URL.Query().Get("user"),
		"update": r.URL.Query().Get("update"),
	}

	response := searchLogic(request)
	bytes, err := json.Marshal(response)
	if err != nil {
		errorResponse(err, w)
		return
	}
	successResponse(bytes, w)
}

func registtaskController(w http.ResponseWriter, r *http.Request) {
	log.Info(r.URL.Path)
	decoder := json.NewDecoder(r.Body)
	var request task
	err := decoder.Decode(&request)
	if err != nil {
		errorResponse(err, w)
		return
	}

	response := registLogic(request)
	bytes, err := json.Marshal(response)
	if err != nil {
		errorResponse(err, w)
		return
	}
	successResponse(bytes, w)
}

func updatetaskController(w http.ResponseWriter, r *http.Request) {
	log.Info(r.URL.Path)
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	var request task
	err := decoder.Decode(&request)
	if err != nil {
		errorResponse(err, w)
		return
	}

	request.Taskid = vars["id"]

	var response updateResponse
	if request.Status != "" {
		response = updatestatusLogic(request)
	} else {
		response = updatedescLogic(request)
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		errorResponse(err, w)
		return
	}
	successResponse(bytes, w)
}

func alluserController(w http.ResponseWriter, r *http.Request) {
	log.Info(r.URL.Path)
	response := alluserLogic()

	bytes, err := json.Marshal(response)
	if err != nil {
		errorResponse(err, w)
		return
	}
	successResponse(bytes, w)
}

func successResponse(res []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(res))
	log.Info("success")
}

func errorResponse(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, error_response, err.Error())
	log.Warn(err)
}
