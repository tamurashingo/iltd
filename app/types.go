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

type task struct {
	Taskid   string `json:"id"`
	Tasktype string `json:"type"`
	Taskname string `json:"name"`
	Plan     string `json:"plan"`
	Result   string `json:"result"`
	Unit     string `json:"unit"`
	Status   string `json:"status"`
	Due      string `json:"due"`
	Person   string `json:"person"`
}

type showTaskResponse struct {
	Result     bool    `json:"result"`
	Todo       *[]task `json:"tasktodo"`
	Inprogress *[]task `json:"taskinprogress"`
	Done       *[]task `json:"taskdone"`
	Archive    *[]task `json:"taskarchive"`
}

type updateResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

type registResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

type user struct {
	Username string `json:"username"`
}

type alluserResponse struct {
	Result bool    `json:"result"`
	Users  *[]user `json:"users"`
}
