package controllers

import (
	"awesomeProject/configs"
	"awesomeProject/models"
	"awesomeProject/responses"
	"awesomeProject/services/mongo_service"
	"awesomeProject/services/redis_service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var validate = validator.New()

func CreateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "app_error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "app_error",
			Data:    &echo.Map{"data": validationErr.Error()},
		})
	}

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	var userCollection = mongo_service.GetCollection(configs.MongoClient, "users")
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "app_error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &echo.Map{"data": result},
	})
}

func GetAUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	var userCollection = mongo_service.GetCollection(configs.MongoClient, "users")
	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "app_error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &echo.Map{"data": user},
	})
}

func GetAllUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	jsonString, redisError := redis_service.Get(c.Path())
	if redisError != nil {
		fmt.Println(redisError.Error())
	} else {
		err := json.Unmarshal([]byte(jsonString), &users)
		if err == nil {
			return c.JSON(http.StatusOK, responses.UserResponse{
				Status:  http.StatusOK,
				Message: "success",
				Data: &echo.Map{
					"data": users,
				},
			})
		}
	}

	var userCollection = mongo_service.GetCollection(configs.MongoClient, "users")
	results, err := userCollection.Find(ctx, bson.D{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "app_error",
			Data:    &echo.Map{"data": err.Error()},
		})
	}

	defer func(results *mongo.Cursor, ctx context.Context) {
		_ = results.Close(ctx)
	}(results, ctx)

	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "app_error",
				Data:    &echo.Map{"data": err.Error()},
			})
		}

		users = append(users, singleUser)
	}

	jsonResults, marshalError := json.Marshal(users)
	if marshalError == nil {
		redis_service.Set(c.Path(), string(jsonResults))
	}

	return c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &echo.Map{
			"data": users,
		},
	})
}
