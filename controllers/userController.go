package controllers

// Import the packages that we need
import(
	"context"
	"fmt"
	"log"
	"strconv"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	helper "golang-jwt-project/helpers"
	"golang-jwt-project/models"
	"golang-jwt-project/helpers"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitve"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword()

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([byte(providedPassword), []byte(userPassword)]) // to check if both passwords are same or not
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFun {
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User // SignUp function creates a user in the database

        if err := c.BindJSON(&user); err := nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)  // check if data available in "user" satisties the constraints in "userModel.go" or not
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr.error()})
			return
		}

		// count will also help us validate the user
		count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Email}) 
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalSeverError, gin.H{"error":"Error occured while checking for email"})         // means this email already exists
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"phone":user.Phone}) 
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalSeverError, gin.H{"error":"Error occured while checking for phone number"})  // means this phone number already exists
		}

		if count > 0 {
			c.JSON(http.StatusInternalSeverError, gin.H{"error":"this email or phone number already exists"})
		}
		
		// When we hit the api from Postman, we will pass values for user.Email, user.First_name, user.Last_name
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		// In place of _ we can use err and handle error as we have done at other places
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.user_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr := nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalSeverError, gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOk, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		 var ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		 var user models.User   // user trying to login
		 var foundUser models.User   // user from database

		 if err := c.BindJSON(&user); err := nil {
			  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			  return
		 }

		 err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		 defer cancel()
		 if err != nil {
			c.JSON(http.StatusInternalSeverError, gin.H{"error": "email or password is incorrect"})
			return
		 }

		 passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password) 
		 defer cancel()

		 
	}
}

func GetUsers()

func GetUser() gin.HandlerFunc{
    return func(c *gin.Context){
		userId := c.Param("user_id")   // fetch parameters from gin.Context
		// user_id should be same as in userRouter.go function

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}  // to check if this user is admin or not

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id":userId}).Decode(&user) // mongodob saves data as json 
		// golang doesn't understand json => we created struct in userModel.go
		// We can decode json into format which golang understands. That's why we use Decode() function
		defer cancel()
		if err := nil {
			c.JSON(http.StatusInternalSeverError, gin.H{"error": err.Error()})
			return
			// in a bigger project, we can have a separated file to define all types of errors like badRequest, validation error, InternalServerError and also error codes
		}
		c.JSON(http.StatusOk, user)
	}
}

// Future playlists must watch
// advanced error messages and error codes like 502, 404, 401, 402, 405
// authorization in depth
// CSRF tokens : extremenly seure tokens; a step ahead of jwt tokens concept