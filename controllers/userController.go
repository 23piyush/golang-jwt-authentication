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
	helper "github.com/23piyush/golang-jwt-project/helpers"
	"github.com/23piyush/golang-jwt-project/models"
	"github.com/23piyush/golang-jwt-project/database"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

// to hash the password before storing it in the database
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword)) // to check if both passwords are same or not
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User // SignUp function creates a user in the database

        if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)  // check if data available in "user" satisties the constraints in "userModel.go" or not
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error()})
			return
		}

		// count will also help us validate the user
		count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Email}) 
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Error occured while checking for email"})         // means this email already exists
		}
    
		password := HashPassword(*user.Password)
		user.Password = &password
        

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone":user.Phone}) 
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Error occured while checking for phone number"})  // means this phone number already exists
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"this email or phone number already exists"})
		}
		
		// When we hit the api from Postman, we will pass values for user.Email, user.First_name, user.Last_name
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		// In place of _ we can use err and handle error as we have done at other places
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		 var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		 var user models.User   // user trying to login
		 var foundUser models.User   // user from database

		 if err := c.BindJSON(&user); err != nil {
			  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			  return
		 }

		 err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		 defer cancel()
		 if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		 }

		 passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password) 
		 defer cancel()
         if passwordIsValid != true {
			 c.JSON(http.StatusInternalServerError, gin.H{"error" : msg})
			 return
		 }

		 if foundUser.Email == nil {
			 c.JSON(http.StatusInternalServerError, gin.H{"error":"user not found"})
		 }
		 token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *&foundUser.User_id)
		 helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		 userCollection.FindOne(ctx, bson.M{"user_id":foundUser.User_id}).Decode(&foundUser)

		 if err != nil {
			 c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			 return
		 }
		 c.JSON(http.StatusOK, foundUser)
	}
}

// This function can be accessed by admins only, as it returns list of users
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})  // not an admin, but trying to access things for admins, thus a bad request
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}
	page, err1 := strconv.Atoi(c.Query("page"))
	if err1 != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(c.Query("startIndex"))

    // mongodb aggregation functions - Match, Project and Group - These are pipeline stage functions
	
	matchStage := bson.D{{"$match", bson.D{{}}}}
	groupStage := bson.D{{"$group", bson.D{
		                                    {"_id", bson.D{{"_id", "null"}}},  // In $group stage, we are grouping documents as per id as unique field
	                                        {"total_count", bson.D{{"$sum", 1}}}, // find total count, we got the count but lost the actual data, thus we need another stage
							                {"data",  bson.D{{"$push", "$$ROOT"}}}}}}
	  // $sum is not a pipeline stage function, it gives here the count of records, which were grouped by id in previous stage, we got actual data here
    projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},}}}
	 // projectStage helps us define which data points should go to the user and which should not go to the user, useful fot front-end, helps us have control ..... much similar to problem solved using GraphQL

       result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
         matchStage, groupStage, projectStage })
	   defer cancel()
	   if err != nil {
		   c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing user items"})
	   }

	   var allUsers []bson.M 
	   if err = result.All(ctx, &allUsers); err!= nil {
		   log.Fatal(err)
	   }
       c.JSON(http.StatusOK, allUsers[0])}}

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
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
			// in a bigger project, we can have a separated file to define all types of errors like badRequest, validation error, InternalServerError and also error codes
		}
		c.JSON(http.StatusOK, user)
	}
}

// Future playlists must watch
// advanced error messages and error codes like 502, 404, 401, 402, 405
// authorization in depth
// CSRF tokens : extremenly seure tokens; a step ahead of jwt tokens concept