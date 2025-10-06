package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"pretest-gis/config"
	"pretest-gis/models"
)
var ctxTimeout = 10 * time.Second

func CreatePlace(c *fiber.Ctx) error {
	var p models.PlacePayload
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}
	if p.Name == "" || p.Latitude == nil || p.Longitude == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name, latitude, longitude required"})
	}

	place := models.Place{
		ID:   primitive.NewObjectID(),
		Name: p.Name,
		Location: models.GeoJSONPoint{
			Type:        "Point",
			Coordinates: []float64{*p.Longitude, *p.Latitude},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	collection := config.DB.Collection("places") // ✅ ambil di sini
	_, err := collection.InsertOne(ctx, place)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db insert error"})
	}
	return c.Status(fiber.StatusCreated).JSON(place)
}

func GetPlaces(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	collection := config.DB.Collection("places") // ✅
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	defer cursor.Close(ctx)

	var places []models.Place
	if err := cursor.All(ctx, &places); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "decode error"})
	}
	return c.JSON(places)
}

func GetPlace(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	collection := config.DB.Collection("places") // ✅
	var place models.Place
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&place)
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	return c.JSON(place)
}

func UpdatePlace(c *fiber.Ctx) error {
    id := c.Params("id")

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
    }

    var payload models.PlacePayload
    if err := c.BodyParser(&payload); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
    }

    update := bson.M{
        "$set": bson.M{
            "name": payload.Name,
            "location": bson.M{
                "type":        "Point",
                "coordinates": []float64{*payload.Longitude, *payload.Latitude},
            },
        },
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := config.DB.Collection("places")
    result, err := collection.UpdateByID(ctx, objID, update)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    if result.MatchedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "place not found"})
    }

    return c.JSON(fiber.Map{"message": "Place updated successfully"})
}


func DeletePlace(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	collection := config.DB.Collection("places") // ✅
	res, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	if res.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}