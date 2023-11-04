package fileservice

import (
	"blackbox-v2/pkg/mongoc"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var FileMetaDataCollection = mongoc.Database.Collection("file_metadata")

type FileMetaData struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	OgFileName  string             `bson:"og_file_name"`
	NewFileName string             `bson:"new_file_name"`
	FilePath    string             `bson:"file_path"`
	UserCID     string             `bson:"user_cid"`
	Extension   string             `bson:"extension"`
	Category    string             `bson:"category"`
	Size        int                `bson:"size"`
	SizeInMB    float64            `bson:"size_in_mb"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

func (fmd *FileMetaData) toBson() bson.M {
	return bson.M{
		"og_file_name":  fmd.OgFileName,
		"new_file_name": fmd.NewFileName,
		"file_path":     fmd.FilePath,
		"user_cid":      fmd.UserCID,
		"extension":     fmd.Extension,
		"category":      fmd.Category,
		"size":          fmd.Size,
		"size_in_mb":    fmd.SizeInMB,
		"created_at":    fmd.CreatedAt,
		"updated_at":    fmd.UpdatedAt,
	}
}

func (fmd *FileMetaData) create() error {
	_, err := FileMetaDataCollection.InsertOne(context.Background(), fmd.toBson())
	return err
}
