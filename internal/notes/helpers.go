package notes

import (
	"blackbox-v2/internal"
	"blackbox-v2/internal/userservice"
	"blackbox-v2/pkg/utils"
	"bufio"
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveFile(file multipart.File, fileHeader *multipart.FileHeader, userCID string) error {
	randomStrng := utils.GenerateRandomString(5)

	fileNameWoExt := strings.Split(fileHeader.Filename, ".")
	fileName, extension := fileNameWoExt[0], fileNameWoExt[1]
	newFileName := fileName + "_" + randomStrng + "." + extension
	user, err := userservice.GetUserByCID(userCID)
	if err != nil {
		return err
	}
	umd, err := userservice.GetUserMetaDataByUserID(user.ID)
	if err != nil {
		return err
	}
	newFilePath := filepath.Join(internal.UploadDir, umd.DirName, "/notes/", newFileName)

	newFile, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	defer newFile.Close()
	writer := io.Writer(newFile)
	reader := bufio.NewReader(file)
	_, err = reader.WriteTo(writer)
	if err != nil {
		return err
	}

	fd := NoteFileMetaData{
		OgFileName: fileHeader.Filename,
		UserCID:    userCID,
		FilePath:   newFilePath,
		Parsed:     false,
	}
	err = fd.create()
	return err
}

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
