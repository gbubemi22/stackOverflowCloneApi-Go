package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Answer struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	Body        *string   `bson:"body" json:"body" binding:"required"`
	User_id     *string   `bson:"user_id" json:"user_id" binding:"required"`
	Question_id *string   `bson:"question_id" json:"user_id" binding:"required"`
	Created_at  time.Time `bson:"created_at" json:"created_at"`
	Updated_at  time.Time `bson:"updated_at" json:"updated_at"`
	Answer_id   string    `bson:"question_id" json:"question_id"`
	Likes       []string  `bson:"likes" json:"likes"`
	Image       *string   `bson:"image" json:"image"`
}
