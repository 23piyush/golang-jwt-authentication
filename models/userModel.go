package models

import(
"time"
"go.mongodb.org/mongo-driver/bson/primitive"
)

    // will act as middle-layer between golang-program and mongoDB-database
    // database understands json and golang doesn't understands, may ne string 
	// we need a layer for inter-conversion
	// In postman or otherwise : "First_name", but in database : "first_name"
	// min 2 characters and max 100 characters
	// validate:"required, min=2, max=100" will give error because you can't have spaces
type User struct{ 
	ID             primitive.ObjectID    `bson:"_id"`
	First_name     *string               `json:"first_name" validate:"required,min=2,max=100"`    
	Last_name      *string                `json:"last_name" validate:"required,min=2,max=100"`
	Password       *string                `json:"Password" validate:"required,min=6"`
	Email          *string                `json:"email" validate:"email,required"`
	Phone          *string                `json:"phone" validate:"required"`
	Token          *string                `json:"token"`
	User_type      *string                `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Refresh_token  *string                `json:"refresh_token"`
	Created_at     time.Time             `json:"created_at"`
	Updated_at     time.Time             `json:"updated_at"`
	User_id        string                `json:"user_id"`
}
