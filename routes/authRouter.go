package routes

import(
	controller "github.com/23piyush/golang-jwt-project/controllers"  // import controllers package
	"github.com/gin-gonic/gin"  // to handle routes
	"fmt"
)

   // gin.Engine represents the main router of the application.
func AuthRoutes(incomingRoutes *gin.Engine){
	// define the routes and the associated function from the controller package
	fmt.Print("here1");
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
	// controller.Signup and controller.Login are called directly with () after their names. This is because the functions are executed immediately and their return values are used as the route handlers
	// Signup and Login are not protected routes because user doesn't have token till now 
	// controller functions defined in the golang-jwt-project/controllers package handle the logic for user signup and login
}
