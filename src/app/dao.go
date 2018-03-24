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
	"container/list"

	log "github.com/Sirupsen/logrus"
)

const (
	// 指定したステータスのタスクを取得する.
	// 追加で以下の条件を指定できる:
	//
	// * `person`. 担当者
	// * `status_update_date`. ステータスを変更した日
	selectSQL = `
 select
   id,
   type,
   name,
   plan,
   result,
   unit,
   due,
   person
 from
   t_task
 where
   task_status = ?
 and
   (? = '' or person = ?)
 and
   (? = '' or status_update_date >= ?)
 order by
   id desc
`
	// タスクを登録する.
	registTaskSQL = `
 insert into
   t_task
 (
   type,
   name,
   plan,
   result,
   unit,
   due,
   person,
   task_status,
   create_date,
   status_update_date,
   desc_update_date
 )
 values (
   ?,
   ?,
   ?,
   ?,
   ?,
   ?,
   ?,
   'task-todo',
   now(),
   now(),
   now()
 )
`
	// ステータスを更新する
	updateStatusSQL = `
 update
   t_task
 set
   task_status = ?,
   status_update_date = now()
 where
   id = ?
`
	// 詳細情報を更新する
	updateDetailSQL = `
 update
   t_task
 set
   type = ?,
   name = ?,
   plan = ?,
   result = ?,
   unit = ?,
   due = ?,
   person = ?,
   desc_update_date = now()
 where
   id = ?
`

	// 全ユーザ名を取得する
	selectAllUserSQL = `
 select
   distinct
   person
 from
   t_task
 where
   task_status in ('task-todo', 'task-inprogress', 'task-done')
`
)

// 指定したステータスのタスクを取得する.
// 追加で以下の条件を指定できる:
//
// * `person`. 担当者
// * `status_update_date`. ステータスを変更した日
func fetchTaskDAO(st string, req map[string]string) (*[]task, error) {
	stmt, err := db.Prepare(selectSQL)

	if err != nil {
		log.WithFields(log.Fields{
			"type":  "database",
			"stage": "prepare",
			"sql":   selectSQL,
		}).Warn(err)
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(
		st,
		req["user"],
		req["user"],
		req["update"],
		req["update"])
	if err != nil {
		log.WithFields(log.Fields{
			"type":  "database",
			"stage": "query",
			"sql":   selectSQL,
		}).Warn(err)
		return nil, err
	}
	defer rows.Close()
	tasks := list.New()

	for rows.Next() {
		var taskid string
		var tasktype string
		var taskname string
		var plan string
		var result string
		var unit string
		var due string
		var person string
		err = rows.Scan(&taskid, &tasktype, &taskname, &plan, &result, &unit, &due, &person)
		if err != nil {
			log.WithFields(log.Fields{
				"type":  "database",
				"stage": "scan",
				"sql":   selectSQL,
			}).Warn(err)
		}
		task := &task{
			Taskid:   taskid,
			Tasktype: tasktype,
			Taskname: taskname,
			Plan:     plan,
			Result:   result,
			Unit:     unit,
			Due:      due,
			Person:   person,
		}
		tasks.PushBack(*task)
	}

	allTask := make([]task, tasks.Len())
	idx := 0
	for v := tasks.Front(); v != nil; v = v.Next() {
		allTask[idx] = v.Value.(task)
		idx = idx + 1
	}

	return &allTask, nil
}

func registTaskDAO(taskinfo task) error {
	stmt, err := db.Prepare(registTaskSQL)

	if err != nil {
		log.WithFields(log.Fields{
			"type":  "database",
			"stage": "prepare",
			"sql":   registTaskSQL,
		}).Warn(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(
		taskinfo.Tasktype,
		taskinfo.Taskname,
		taskinfo.Plan,
		taskinfo.Result,
		taskinfo.Unit,
		taskinfo.Due,
		taskinfo.Person,
	)

	if err != nil {
		log.WithFields(log.Fields{
			"type":  "database",
			"stage": "exec",
			"sql":   registTaskSQL,
		}).Warn(err)
		return err
	}

	return nil
}

func updatestatusTaskDAO(taskid string, taskst string) error {
	stmt, err := db.Prepare(updateStatusSQL)

	if err != nil {
		log.WithFields(log.Fields{
			"type":  "database",
			"stage": "prepare",
			"sql":   updateStatusSQL,
		}).Warn(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(taskst, taskid)

	if err != nil {
		log.WithFields(log.Fields{
			"type":  "database",
			"stage": "exec",
			"sql":   updateStatusSQL,
		}).Warn(err)
		return err
	}

	return nil
}

func updatedescTaskDAO(taskinfo task) error {
	stmt, err := db.Prepare(updateDetailSQL)

	if err != nil {
		log.WithFields(log.Fields{
			"type":  "database",
			"stage": "prepare",
			"sql":   updateDetailSQL,
		}).Warn(err)
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(
		taskinfo.Tasktype,
		taskinfo.Taskname,
		taskinfo.Plan,
		taskinfo.Result,
		taskinfo.Unit,
		taskinfo.Due,
		taskinfo.Person,
		taskinfo.Taskid,
	)

	if err != nil {
		log.WithFields(log.Fields{
			"type":  "database",
			"stage": "exec",
			"sql":   updateDetailSQL,
		}).Warn(err)
		return err
	}

	return nil
}

func allUserDAO() (*[]user, error) {
	stmt, err := db.Prepare(selectAllUserSQL)

	if err != nil {
		log.WithFields(log.Fields{
			"type":  "database",
			"stage": "prepare",
			"sql":   selectAllUserSQL,
		}).Warn(err)
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.WithFields(log.Fields{
			"type":  "database",
			"stage": "prepare",
			"sql":   selectAllUserSQL,
		}).Warn(err)
		return nil, err
	}
	defer rows.Close()
	users := list.New()

	for rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			log.WithFields(log.Fields{
				"type":  "database",
				"stage": "scan",
				"sql":   selectAllUserSQL,
			}).Warn(err)
		}
		user := &user{
			Username: username,
		}
		users.PushBack(*user)
	}

	allUsers := make([]user, users.Len())
	idx := 0
	for v := users.Front(); v != nil; v = v.Next() {
		allUsers[idx] = v.Value.(user)
		idx = idx + 1
	}

	return &allUsers, nil
}
