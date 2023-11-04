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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func ParseNotes() error {
	var noteMds []NoteFileMetaData
	var err error
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})

	curr, err := NotesFileMetaDataCollection.Find(context.TODO(), bson.D{{"parsed", false}}, opts)
	if err != nil {
		return err
	}
	var results []bson.M
	if err = curr.All(context.TODO(), &results); err != nil {
		return err
	}
	for _, result := range results {
		noteMds = append(noteMds, NoteFileMetaData{
			ID:         result["_id"].(primitive.ObjectID),
			OgFileName: result["og_file_name"].(string),
			FilePath:   result["file_path"].(string),
			UserCID:    result["user_cid"].(string),
			NoteCID:    result["note_cid"].(string),
			Parsed:     result["parsed"].(bool),
		})
	}
	for _, noteMd := range noteMds {
		err = parseNoteFile(noteMd)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseNoteFile(nmd NoteFileMetaData) error {
	file, err := os.Open(nmd.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	note := Note{
		Heading:   lines[0],
		Text:      strings.Join(lines[1:], "\n"),
		UserCID:   nmd.UserCID,
		CID:       nmd.NoteCID,
		CreatedAt: nmd.CreatedAt,
	}
	err = note.create()
	if err != nil {
		return err
	}
	nmd.Parsed = true
	opts := options.FindOneAndUpdate()
	filter := bson.D{{"_id", nmd.ID}}
	update := bson.D{{"$set", bson.D{{"parsed", true}}}}
	var updatedDoc bson.M
	err = NotesFileMetaDataCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedDoc)

	return err
}
