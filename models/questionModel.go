package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Question struct {
	ID          primitive.ObjectID `bson:"id"`
	Title       *string            `json:"title" validate:"required,max=50"`
	Body        *string            `json:"body" validate:"required,min=3"`
	Likes       []string           `json:"likes"`
	Image       *string            `json:"image" `
	Question_id string             `json:"question_id"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	User_id     *string            `json:"user_id" validate:"required`
}
