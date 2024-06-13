# Ghttp

Ghttp stands for Good Http.

Is a VERY VERY simple http framework inspired in Echo.
Echo framework return error on its handlers, and thats something that I believe to be really clever. So thats what I did here

The reason why I created this project was to learn something new and to also simplify the development of small projects that a simple http package would do the trick... Now with the http package from 1.22, for me all it was missing was the possibility to return errors in the handlers, so I just this one.

Example:
``` golang
  package main

  import (
    "github.com/victorguidi/ghttp/ghttp"
  )

  func main() {
    server := ghttp.New().CORS()
    server.GET("/", get)
    server.Start(":5000")
  }

  func get(c ghttp.Context) error {
    type Response struct {
      Name string `json:"name"`
    }
    var resp Response
    resp.Name = "John Wick"

    return c.JSON(resp)
  }
```
