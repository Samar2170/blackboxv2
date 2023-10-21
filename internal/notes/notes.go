package notes

import (
	"blackbox-v2/pkg/mongoc"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var NotesCollection = mongoc.Database.Collection("notes")
var NotesFileMetaDataCollection = mongoc.Database.Collection("notes_file_metadata")

type NoteFileMetaData struct {
	ID         uint
	OgFileName string
	FilePath   string
	UserCID    string
	NoteCID    string
	Parsed     bool

	CreatedAt time.Time
	UpdatedAt time.Time
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
	ID      uint
	Heading string
	Text    string
	UserCID string
	CID     string

	CreatedAt time.Time
	UpdatedAt time.Time
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
