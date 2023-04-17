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

var answerCollection *mongo.Collection = database.OpenCollection(database.Client, "answer")
//var validate = validator.New()

func CreatAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var question models.Question
		var user models.User
		var answer models.Answer

		if err := c.BindJSON(&answer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validationErr := validate.Struct(answer)
		// if validationErr != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		// 	return
		// }
		err := userCollection.FindOne(ctx, bson.M{"user_id": answer.User_id}).Decode(&user)
		defer cancel()

		if err != nil {
			msg := fmt.Sprintf("user was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// err := qestionCollection.FindOne(ctx, bson.M{"questionId": answer.QuestionId}).Decode(&question)
		// defer cancel()

		// if err != nil {
		// 	msg := fmt.Sprintf("qestion was not found")
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		// 	return
		// }

		answer.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		answer.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		answer.ID = primitive.NewObjectID()
		answer.Answer_id = answer.ID.Hex()

		result, insertErr := answerCollection.InsertOne(ctx, question)

		if insertErr != nil {
			msg := fmt.Sprintf("Qestion item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}


func GetOneAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		answerId := c.Param("answer_id")
		var answer models.Answer

		err := answerCollection.FindOne(ctx, bson.M{"answer_id": answerId}).Decode(&answer)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching this Answer"})
		}
		c.JSON(http.StatusOK, answer)
	}
}


func GetAllQAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := answerCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fatching All Answers"})
		}

		var allAnswers []bson.M
		if err = result.All(ctx, &allAnswers); err != nil {

			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allAnswers)
		return
	}
}


func UpdateAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var question models.Question
		var user models.User
		var answer models.Answer

		answerId := c.Param("answer_id")

		if err := c.BindJSON(&question); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D

		if answer.Body != nil {
			updateObj = append(updateObj, bson.E{"title", question.Title})
		}

		


		if answer.User_id != nil {
			err := userCollection.FindOne(ctx, bson.M{"user_id": question.User_id}).Decode(&user)
			defer cancel()
			
			if err != nil {
				msg := fmt.Sprintf("message:User was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}

		}

		if answer.Question_id != nil {
			err := qestionCollection.FindOne(ctx, bson.M{"question_id": answer.Question_id}).Decode(&question)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("message:User was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}

		}

		answer.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", answer.Updated_at})

		upsert := true
		filter := bson.M{"answer_id": answerId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := answerCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprint(" answer update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}

func UpdateAnswerLikes() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")
		answer_id := c.Param("answer_id ")
		var answer models.Answer
		
		
		// Retrieve question document using question_id
		filter := bson.M{"id": answer_id}
		update := bson.M{"$addToSet": bson.M{"likes": user_id}}
		
		err := qestionCollection.FindOneAndUpdate(
			context.Background(),
			filter,
			update,
			options.FindOneAndUpdate().SetReturnDocument(options.After),
		).Decode(&answer)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if containsA(answer.Likes, user_id) {
			update = bson.M{"$pull": bson.M{"likes": user_id}}
			_, err = answerCollection.UpdateOne(
				context.Background(),
				filter,
				update,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			answer.Likes = removeA(answer.Likes, user_id)
		} else {
			answer.Likes = append(answer.Likes, user_id)
		}

		update = bson.M{"$set": bson.M{"likes": answer.Likes, "updated_at": time.Now()}}
		_, err = answerCollection.UpdateOne(
			context.Background(),
			filter,
			update,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, answer)
	}
}



func containsA(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}

	return false
}

func removeA(arr []string, val string) []string {
	for i, item := range arr {
		if item == val {
			return append(arr[:i], arr[i+1:]...)
		}
	}

	return arr
}