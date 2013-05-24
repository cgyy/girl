package main

import (
    "github.com/cgyy/girl"
)

func Index(c *girl.Context) girl.View {
	return c.RenderText("hello world")
}

func main() {
    girl.Get("/", Index)
    girl.Run(":9999")
}
