package main

import (
	"fmt"
	"go18_projects/book/v3/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	server_engine := gin.Default()

	handlers.Book.Registry(server_engine)

	if err := server_engine.Run(":8080"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
