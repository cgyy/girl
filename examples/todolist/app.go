package main

import (
    "github.com/cgyy/girl"
    "fmt"
    "github.com/cgyy/girl/examples/todolist/model"
)

func Auth(c *girl.Context) girl.View {
    return nil
}

// GET /
func Index(c *girl.Context) girl.View {
    return c.Render("index", nil)
}

// GET /tasks
func GetTasks(c *girl.Context) girl.View {
    tasks := model.FindTasks()
    return c.RenderJSON(tasks)
}

// GET /tasks/:id
func GetTask(c *girl.Context) girl.View {
    return nil
}

// POST /tasks
func PostTask(c *girl.Context) girl.View {
    return nil
}

// PUT /tasks/:id
func PutTask(c *girl.Context) girl.View {
    return nil
}

// PUT /tasks/:id
func PutTask(c *girl.Context) girl.View {
    return nil
}



func main() {
    intro := `
******************************************************************************
a simple todo list
to run this: 
    1.make sure you have a mysql installed, and run "source todolist.sql"
    2.go build app.go && ./app

******************************************************************************

    `
    fmt.Println(intro)

    app := girl.New()
    app.Get("/", Index)
    app.Get("/tasks", GetTasks)

    model.Init()
    app.Run(":9999")
}
