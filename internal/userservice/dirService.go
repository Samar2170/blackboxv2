package userservice

import (
	"blackbox-v2/internal"
	"blackbox-v2/pkg/utils"
	"os"
)

func CreateDirForUser(user *User) error {
	username := user.Username
	dirEntries, err := os.ReadDir(internal.UploadDir)
	if err != nil {
		return err
	}
	var dirNames []string
	for _, dir := range dirEntries {
		dirNames = append(dirNames, dir.Name())
	}
	for i := 0; i < 10; i++ {
		if utils.ArrayContains(dirNames, username) {
			username = username + utils.GenerateRandomString(5)
		} else {
			os.Mkdir(internal.UploadDir+"/"+username, 0777)
			os.Mkdir(internal.UploadDir+"/"+username+"/notes", 0777)
			umd, err := GetUserMetaDataByUserID(user.ID)
			if err != nil {
				return err
			}
			umd.DirCreated = true
			umd.DirName = username
			err = umd.Update()
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}
