package main

import (
	"fmt"
	"github.com/cgyy/girl"
	"github.com/cgyy/girl/examples/todolist/model"
	"strconv"
)

// TODO not used
func Auth(c *girl.Context) girl.View {
	ck, _ := c.Request.Cookie("user_id")
	if ck == nil {
		return c.Redirect("/login")
	}
	return nil
}

// GET /
func index(c *girl.Context) girl.View {
	return c.Render("index", model.User{Name: model.GetUser()})
}

// GET /tasks
func getTasks(c *girl.Context) girl.View {
	tasks := model.FindTasks()
	return c.RenderJSON(tasks)
}

// GET /tasks/:id
func getTask(c *girl.Context) girl.View {
	id, _ := strconv.Atoi(c.GetParam("id"))

	task := model.FindTask(id)
	return c.Render("task", task)
}

// POST /tasks
func addTask(c *girl.Context) girl.View {
	task := model.Task{
		UserId:  c.GetNumParam("user_id"),
		Title:   c.GetParam("title"),
		Content: c.GetParam("content"),
		Active:  true,
	}

	task.Save()
	return c.RenderJSON(task)
}

// PUT /tasks/:id
func updateTask(c *girl.Context) girl.View {
	active := true
	if c.GetParam("active") != "0" {
		active = false
	}
	task := model.Task{
		Id:      c.GetNumParam("id"),
		UserId:  c.GetNumParam("user_id"),
		Title:   c.GetParam("title"),
		Content: c.GetParam("content"),
		Active:  active,
	}

	task.Save()
	return c.RenderJSON(task)
}

// DELETE /tasks/:id
func deleteTask(c *girl.Context) girl.View {
	model.DeleteTask(c.GetNumParam("id"))
	return c.RenderText("success")
}

/**
******************************************************************************
a simple todo list
to run this:
1.make sure you have a mysql installed, and run sql:

CREATE TABLE `tasks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `title` varchar(255) DEFAULT NULL,
  `content` varchar(255) DEFAULT NULL,
  `active` tinyint(1) DEFAULT '1',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8

2. export dbn= "user:password@host/database?charset=utf8"
  else default dbn is "root:123456@/todolist?charset=utf8"

3.go build app.go && ./app

******************************************************************************
**/

func main() {
	fmt.Println("run todo list")

	app := girl.New()

	//app.Before("/.*", Auth) TODO add filter

	app.Get("/", index)
	app.Get("/tasks", getTasks)
	app.Get("/tasks/:id", getTask)
	app.Post("/tasks", addTask)
	app.Put("/tasks/:id", updateTask)
	app.Delete("/tasks/:id", deleteTask)

	model.Init()
	app.Run(":9999")
}
