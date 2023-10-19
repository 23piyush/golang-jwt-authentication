package main

import(
	routes "github.com/23piyush/golang-jwt-project/routes"   // import routes package to handle authentication and user routes
	"os"     // os package is needed for environment variable access
	"log"
	"github.com/gin-gonic/gin"  // gin package is needed for creating the server and handling requests
	"github.com/joho/godotenv"
)
 // main function sets up the server
func main(){
	err := godotenv.Load(".env")
	
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
        port = "8080"
	}

	//gin creates a router for us, similar to gorilla-MUX router, fiber package
	router := gin.New() // creates a new gin router
	router.Use(gin.Logger())  // middleware is added to log incoming requests

	// define additional routes
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// In router handler functions, you don't need to pass req i.e. (w) and res; you can access req and res from gin.Context
	// A route is defined for `/api-1` using `router.GET`
	// When a GET request is made to this route, the provided anonymous function is executed. It responds with a JSON payload containing a success message.
	router.GET("/api-1", func(c *gin.Context) {
        c.JSON(200, gin.H{"success":"Access granter for api-1"})  // we don't need to use w.setHeader=100, gin will do this for us
	})

	// To learn how to create your own router library, don't use any fancy tool like gin, fiber
	// Stick to basic packages like GORILLA MUX 
	// tools save time, but learning won't be there 
	// framework changes with time, thus learn the concept behind them, modifying them and creating your own

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success":"Access granted for api-2"})
	})

	router.Run(":" + port)  // starts the server
}