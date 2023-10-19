package helper

import(
	"context"
	"fmt" // to print out stuffs
	"log"  // to print out errors
	"os"
	"time"
	"golang-jwt-project/database"
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

func GenerateAllTokens (email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string) {
	claims := $SignedDetails{
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
          ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).unix(),
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
