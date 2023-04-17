package controllers

import (
	"context"
	"fmt"
	database "github.com/gbubemi22/go-stackOverflow/database"
	"github.com/gbubemi22/go-stackOverflow/models"
	"log"

	"net/http"
	//"strconv"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var qestionCollection *mongo.Collection = database.OpenCollection(database.Client, "question")
//var validate = *validator.New()

func CreatQestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var question models.Question
		var user models.User

		if err := c.BindJSON(&question); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validationErr := validate.Struct(question)
		// if validationErr != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		// 	return
		// }
		err := userCollection.FindOne(ctx, bson.M{"user_id": question.User_id}).Decode(&user)
		defer cancel()

		if err != nil {
			msg := fmt.Sprintf("user was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		question.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		question.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		question.ID = primitive.NewObjectID()
		question.Question_id = question.ID.Hex()

		

		result, insertErr := qestionCollection.InsertOne(ctx, question)

		if insertErr != nil {
			msg := fmt.Sprintf("Qestion item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func GetOneQestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		questionId := c.Param("question_id")
		var question models.Question

		err := qestionCollection.FindOne(ctx, bson.M{"question_id": questionId}).Decode(&question)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching this qestion"})
		}
		c.JSON(http.StatusOK, question)
	}
}

func GetAllQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := qestionCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fatching All Qestions"})
		}

		var allQestions []bson.M
		if err = result.All(ctx, &allQestions); err != nil {

			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allQestions)
		return
	}
}

func UpdateQestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var question models.Question
		var user models.User

		qestionId := c.Param("qestion_id")

		if err := c.BindJSON(&question); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if question.Title != nil {
			updateObj = append(updateObj, bson.E{"title", question.Title})
		}

		if question.Body != nil {
			updateObj = append(updateObj, bson.E{"body", question.Body})
		}


		if question.User_id != nil {
			err := userCollection.FindOne(ctx, bson.M{"user_id": question.User_id}).Decode(&user)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("message:User was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}

		}

		question.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", question.Updated_at})

		upsert := true
		filter := bson.M{"question_id": qestionId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := qestionCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprint(" question update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}



func UpdateLikes() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")
		question_id := c.Param("question_id")
		var question models.Question

		// Retrieve question document using question_id
		filter := bson.M{"id": question_id}
		update := bson.M{"$addToSet": bson.M{"likes": user_id}}
		
		err := qestionCollection.FindOneAndUpdate(
			context.Background(),
			filter,
			update,
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&question)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if contains(question.Likes, user_id) {
			update = bson.M{"$pull": bson.M{"likes": user_id}}
			_, err = qestionCollection.UpdateOne(
				context.Background(),
				filter,
				update,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			question.Likes = remove(question.Likes, user_id)
		} else {
			question.Likes = append(question.Likes, user_id)
		}

		update = bson.M{"$set": bson.M{"likes": question.Likes, "updated_at": time.Now()}}
		_, err = qestionCollection.UpdateOne(
			context.Background(),
			filter,
			update,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, question)
	}
}

func contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}

	return false
}

func remove(arr []string, val string) []string {
	for i, item := range arr {
		if item == val {
			return append(arr[:i], arr[i+1:]...)
		}
	}

	return arr
}