package models

import(
"time"
"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{ 
	// will act as middle-layer between golang-program and mongoDB-database
    // database understands json and golang doesn't understands, may ne string 
	// we need a layer for inter-conversion
	ID             primitive.objectID    `bson:"_id`
	First_name     *string               `json:"first_name" validate:"required, min=2, max=100"`
	// In postman or otherwise : "First_name", but in database : "first_name"
	// min 2 characters and max 100 characters
	Last_name      *string                `json:"first_name" validate:"required, min=2, max=100"`
	Password       *string                `json:"Password" validate:"required, min=6"`
	Email          *string                `json:"email" validate:"email, required"`
	Phone          *string                `json:"phone" validate:"requred"`
	Token          *string                `json:"token"`
	User_type      *string                `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	Refresh_token  *string                `json:"refresh_token"`
	Created_at     *time.time             `json:"created_at"`
	Updated_at     *time.time             `json:"updated_at"`
	User_id        *string                `json:"user_id"`
}
