package routes

import(
	controller "golang-jwt-project/controllers"  // import controllers package
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
	// Signup and Login are not protected routes because user doesn't have token till now 
}