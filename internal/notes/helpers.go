package notes

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListNote struct {
	ID        primitive.ObjectID `bson:"_id"`
	Heading   string             `bson:"heading"`
	UserCID   string             `bson:"user_cid"`
	CID       string             `bson:"cid"`
	CreatedAt time.Time          `bson:"created_at"`
}

func GetNoteMetaDataByUser(userCID string) ([]ListNote, error) {
	var notes []ListNote
	c, err := NotesCollection.Find(context.Background(), bson.D{{"user_cid", userCID}})
	if err != nil {
		return notes, err
	}
	err = c.All(context.Background(), &notes)
	return notes, err
}

func GetNoteByNoteID(id string) (Note, error) {
	var note Note
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return note, err
	}

	err = NotesCollection.FindOne(context.Background(), bson.D{{"_id", objId}}).Decode(&note)
	return note, err
}
