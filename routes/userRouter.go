package routes

import(
	controller "github.com/23piyush/golang-jwt-project/controllers"
	"github.com/23piyush/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())  // middleware is used to ensure that both these routes are protected routes. User has token after login and hence he shouldn't be able to access these protected routes
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}