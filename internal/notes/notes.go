package notes

import (
	"blackbox-v2/pkg/db"
	"blackbox-v2/pkg/mongoc"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	db.DB.AutoMigrate(&Note{}, &NoteFileMetaData{})
}

var NotesCollection = mongoc.Database.Collection("notes")
var NotesFileMetaDataCollection = mongoc.Database.Collection("notes_file_metadata")

type NoteFileMetaData struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	OgFileName string             `bson:"og_file_name"`
	FilePath   string             `bson:"file_path"`
	UserCID    string             `bson:"user_cid"`
	NoteCID    string             `bson:"note_cid"`
	Parsed     bool               `bson:"parsed"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (n *NoteFileMetaData) toBson() bson.M {
	return bson.M{
		"og_file_name": n.OgFileName,
		"file_path":    n.FilePath,
		"user_cid":     n.UserCID,
		"note_cid":     n.NoteCID,
		"parsed":       n.Parsed,
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
	}
}

func (n *NoteFileMetaData) create() error {
	_, err := NotesFileMetaDataCollection.InsertOne(context.Background(), n.toBson())
	return err
}

type Note struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Heading string             `bson:"heading"`
	Text    string             `bson:"text"`
	UserCID string             `bson:"user_cid"`
	CID     string             `bson:"cid"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (n *Note) toBson() bson.M {
	return bson.M{
		"heading":    n.Heading,
		"text":       n.Text,
		"user_cid":   n.UserCID,
		"cid":        n.CID,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
}

func (n *Note) create() error {
	_, err := NotesCollection.InsertOne(context.Background(), n.toBson())
	return err
}
