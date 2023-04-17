package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       *string            `bson:"title,omitempty" json:"title,omitempty" binding:"required"`
	Body        *string            `bson:"body,omitempty" json:"body,omitempty" binding:"required"`
	User_id     *string            `bson:"user_id,omitempty" json:"user_id,omitempty" binding:"required"`
	Created_at  time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at  time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	Question_id string             `bson:"question_id,omitempty" json:"question_id,omitempty"`
	Likes       []string           `bson:"likes,omitempty" json:"likes,omitempty"`
}
