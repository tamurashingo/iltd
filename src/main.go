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
package main

import (
	"os"

	"github.com/codegangsta/cli"

	log "github.com/Sirupsen/logrus"

	"./app"
)

func main() {
	cmd := cli.NewApp()
	cmd.Version = "0.0.1"
	cmd.Name = "iltd"
	cmd.Usage = "IL task management server"
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "set configuration file",
			Value: "config.json",
		},
	}

	cmd.Action = func(c *cli.Context) {
		filename := c.String("config")
		iltd, err := app.IltdNew(filename)

		if err != nil {
			log.Fatal("initialize error:", err)
			return
		}

		err = iltd.IltdRun()
		if err != nil {
			log.Fatal("iltd error:", err)
		}
	}

	cmd.Run(os.Args)
}
