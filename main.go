package main

import(
	routes "golang-jwt-project/routes"   // import routes package
	"os"
	"github.com/gin-gonic/gin"
)

func main(){
	port = os.Getenv("PORT")

	if port == "" {
        port = "8080"
	}

	//gin creates a router for us, similar to gorilla-MUX router, fiber package
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// In router handler functions, you don't need to pass req i.e. (w) and res; you can access req and res from
	// gin.Context
	router.GET("/api-1", func(c *gin.Context)) {
        c.JSON(200, gin.H{"success":"Access granter for api-1"})  // we don't need to use w.setHeader=100, gin will do this for us
	}

	// To learn how to create your own router library, don't use any fancy tool like gin, fiber
	// Stick to basic packages like GORILLA MUX 
	// tools save time, but learning won't be there 
	// framework changes with time, thus learn the concept behind them, modifying them and creating your own

	router.GET("/api-1", func(c *gin.Context)) {
		c.JSON(200, gin.H{"success":"Access granted for api-1"})
	}

	router.Run(":" + port)
}