package fileservice

import (
	"blackbox-v2/internal"
	"blackbox-v2/internal/userservice"
	"blackbox-v2/pkg/logging"
	"blackbox-v2/pkg/utils"
	"bufio"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveFile(file multipart.File, fileHeader *multipart.FileHeader, userCID string) error {
	randomString := utils.GenerateRandomString(5)
	fileNameWoExt := strings.Split(fileHeader.Filename, ".")
	fileName, extension := fileNameWoExt[0], fileNameWoExt[1]
	newFileName := fileName + "_" + randomString + "." + extension

	if _, ok := ValidExtensions[extension]; !ok {
		return errors.New("unallowed file extension")
	}

	if _, ok := CategoryExtensions[extension]; !ok {
		logging.BlackboxCLILogger.Println("my bad, we havent mapped this extension ")
		logging.BlackboxCLILogger.Printf("%s faced with filename %s for extension %s at %s",
			userCID, fileName, extension, time.Now().String())
	}
	category := CategoryExtensions[extension]
	user, err := userservice.GetUserByCID(userCID)
	if err != nil {
		return err
	}
	umd, err := userservice.GetUserMetaDataByUserID(user.ID)
	if err != nil {
		return err
	}
	newFilePath := filepath.Join(internal.UploadDir, umd.DirName, newFileName)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	defer newFile.Close()
	reader := bufio.NewReader(file)
	writer := io.Writer(newFile)
	_, err = reader.WriteTo(writer)
	if err != nil {
		return err
	}
	fileSize := reader.Size()
	fileSizeInMB := utils.ConvertFileSize(float64(fileSize), "bytes", "mb")
	fmd := FileMetaData{
		OgFileName:  fileName,
		NewFileName: newFileName,
		FilePath:    newFilePath,
		UserCID:     userCID,
		Extension:   extension,
		Category:    category,
		Size:        fileSize,
		SizeInMB:    fileSizeInMB,
		CreatedAt:   time.Now(),
	}
	err = fmd.create()
	if err != nil {
		return err
	}
	return nil
}

func GetFilesByUser(userCID string) ([]FileMetaData, error) {
	var fmds []FileMetaData
	c, err := FileMetaDataCollection.Find(
		context.Background(),
		bson.D{{"user_cid", userCID}},
	)
	if err != nil {
		return fmds, err
	}
	err = c.All(context.TODO(), &fmds)
	if err != nil {
		return fmds, err
	}
	return fmds, nil
}

func GetFileByID(id string) ([]byte, FileMetaData, error) {
	var fmd FileMetaData
	var data []byte

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return data, fmd, err
	}
	err = FileMetaDataCollection.FindOne(
		context.TODO(),
		bson.D{{"_id", objId}}).Decode(&fmd)

	data, err = ioutil.ReadFile(fmd.FilePath)
	return data, fmd, err
}
