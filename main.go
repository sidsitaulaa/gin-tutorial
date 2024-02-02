package main

import (
	"encoding/xml"

	"github.com/gin-gonic/gin"
)

type Person struct {
	XMLName  xml.Name `xml:"user"`
	UserName string   `xml:"uname,attr"`
	Age      int32    `xml:"age,attr"`
}

func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello world",
	})
}

func UserHandler(c *gin.Context) {
	name := c.Params.ByName("name")
	c.JSON(200, gin.H{
		"message": "Hello " + name,
	})
}

func UserHandlerXML(c *gin.Context) {
	c.XML(200, Person{
		UserName: "sid",
		Age:      22,
	})
}

func main() {
	router := gin.Default()

	router.GET("/", IndexHandler)
	router.GET("/:name", UserHandler)
	router.GET("/user", UserHandlerXML)

	router.Run(":5000")
}
