package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/KushalSr/calories-tracker-app/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

var entryCollection *mongo.Collection = openCollection(Client, "calories")

//// Add Entry

func AddEntry(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry models.Entry
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validatationError := validate.Struct(entry)

	if validatationError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validatationError.Error()})
		fmt.Println(validatationError)
		return
	}
	entry.ID = primitive.NewObjectID()
	result, insertErr := entryCollection.InsertOne(ctx, entry)
	if insertErr != nil {
		msg := fmt.Sprint("Order Item was not Created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, result)
}

//// Get Entries

func GetEntries(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(err)
	c.JSON(http.StatusOK, entries)
}

//// Get One Entry

func GetEntryByID(c *gin.Context) {
	EntryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(EntryID)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var entry bson.M
	if err := entryCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(entry)
	c.JSON(http.StatusOK, entry)
}

// Get Ingredient

func GetIngredientById(c *gin.Context) {
	ingredient := c.Params.ByName("id")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{"ingredients": ingredient})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()
	fmt.Println(entries)

	c.JSON(http.StatusOK, entries)
}

//// Update Entry

func UpdateEntry(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var entry models.Entry

	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validatationError := validate.Struct(entry)
	if validatationError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validatationError.Error()})
		fmt.Println(validatationError)
		return
	}
	result, err := entryCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docID},
		bson.M{
			"dish":        entry.Dish,
			"fat":         entry.Fat,
			"ingredients": entry.Ingredients,
			"calories":    entry.Calories,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result.ModifiedCount)
}

// Update Ingredient

func UpdateIngredient(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	type Ingredient struct {
		Ingredients *string `json:"ingredients"`
	}
	var ingredient Ingredient

	if err := c.BindJSON(&ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	result, err := entryCollection.UpdateOne(ctx, bson.M{"_id": docID},
		bson.D{{"$set", bson.D{{"ingredients", ingredient.Ingredients}}}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, result.ModifiedCount)
}

/////// Delete Entry

func DeleteEntry(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	result, err := entryCollection.DeleteOne(ctx, bson.M{"_id": docID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)
}
