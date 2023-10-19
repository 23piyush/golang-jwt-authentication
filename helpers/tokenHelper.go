package helper

import(
	"context"
	"fmt" // to print out stuffs
	"log"  // to print out errors
	"os"
	"time"
	"github.com/23piyush/golang-jwt-project/database"
	jwt "github.com/dgrijalva/jwt-go" // In nodejs, you have to do npm for jwt
	// someone has created a golang-driver for jwt for us which we will use
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// jwt token uses a hashing mechanism to take the details that we give it to give us a token
// If we go to jwt.io website, we can decode the token and get all the details back
// a jwt token also consists of a secret key
type SignedDetails struct {
	 Email string
	 First_name string
	 Last_name string
	 Uid string
	 User_type string
	 jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens (email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email : email,
		First_name : firstName,
		Last_name : lastName,
		Uid : uid,
		User_type : userType,
		StandardClaims : jwt.StandardClaims{
			ExpiresAt :time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	 
	// refreshToken is used to get a new token, if our regular token has expired
	refreshClaims := &SignedDetails{
		StandardClaims : jwt.StandardClaims{
          ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	// jwt.SigningMethodHS256 : algorithm to create encrypted token for us
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err 
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token)(interface{}, error) {
             return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid") // someone is sending wrong tokens to access those routes
		msg = err.Error()
		return
	}

	// claims has all the information that the user has
    if claims.ExpiresAt < time.Now().Local().Unix() { // check if the token has expired
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg
}

// Everytime you login, you will get a new token, new refreshed token
func UpdateAllTokens (signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

	upsert := true
	filter := bson.M{"user_id" : userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	 _ , err := userCollection.UpdateOne(
		ctx,  // to update that particular user
		filter, // using userid
		bson.D{
			{"$set", updateObj},   // You will get syntax error without this ","
		},
		&opt,
	 )
	 defer cancel()

	 if err != nil {
		 log.Panic(err)
		 return
	 }

	 return
}
