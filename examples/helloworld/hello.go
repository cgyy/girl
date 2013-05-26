package main

import (
    "github.com/cgyy/girl"
)

func Index(c *girl.Context) girl.View {
	return c.RenderText("hello world")
}

func main() {
    app := girl.New()

    app.Get("/", Index)
    app.Run(":9999")
}
