# girl

girl is a sinatra style web framework 

## Overview


## Installation

Make sure you have the a working Go environment. See the [install instructions](http://golang.org/doc/install.html). 

To install girl, simply run:

    go get github.com/cgyy/girl



## Example
```go
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


```

To run the application, put the code in a file called hello.go and run:

    go run hello.go
    

