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

func searchLogic(req map[string]string) showTaskResponse {
	response := showTaskResponse{}
	tables := map[string](func(t *[]task)){
		"task-todo": func(t *[]task) {
			response.Todo = t
		},
		"task-inprogress": func(t *[]task) {
			response.Inprogress = t
		},
		"task-done": func(t *[]task) {
			response.Done = t
		},
		"task-archive": func(t *[]task) {
			response.Archive = t
		},
	}

	for status, fn := range tables {
		var err error
		result, err := fetchTaskDAO(status, req)
		if err != nil {
			response.Result = false
			return response
		}
		fn(result)
	}

	response.Result = true
	return response
}

func updatestatusLogic(req task) updateResponse {
	err := updatestatusTaskDAO(req.Taskid, req.Status)
	if err != nil {
		return updateResponse{
			Result:  false,
			Message: "タスクの更新に失敗しました",
		}
	}

	return updateResponse{
		Result:  true,
		Message: "タスクの更新に成功しました",
	}
}

func updatedescLogic(req task) updateResponse {
	err := updatedescTaskDAO(req)
	if err != nil {
		return updateResponse{
			Result:  false,
			Message: "タスクの更新に失敗しました",
		}
	}

	return updateResponse{
		Result:  true,
		Message: "タスクの更新に成功しました",
	}
}

func registLogic(req task) registResponse {
	err := registTaskDAO(req)
	if err != nil {
		return registResponse{
			Result:  false,
			Message: "タスク登録に失敗しました",
		}
	}

	return registResponse{
		Result:  true,
		Message: "タスクを登録しました",
	}
}

func alluserLogic() alluserResponse {
	response := alluserResponse{}

	result, err := allUserDAO()
	if err != nil {
		return alluserResponse{
			Result: false,
			Users:  nil,
		}
	}

	response.Result = true
	response.Users = result

	return response
}
