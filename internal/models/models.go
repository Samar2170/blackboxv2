package models

import (
	"errors"
	"reflect"

	"gorm.io/gorm"
)

type DBModel interface {
	Create() error
	Update() error
}

func CreateModelInstance(instance DBModel) error {
	err := instance.Create()
	modelName := reflect.TypeOf(instance).Name()
	switch err {
	case nil:
		return nil
	case gorm.ErrDuplicatedKey:
		return errors.New("[gorm]:ErrDuplicatedKey" + modelName + " already exists")
	case gorm.ErrForeignKeyViolated:
		return errors.New("[gorm]:ErrForeignKeyViolated" + modelName + " missing foreign key")
	default:
		return err
	}
}

func UpdateModelInstance(instance DBModel) error {
	err := instance.Update()
	modelName := reflect.TypeOf(instance).Name()
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return errors.New("[gorm]:ErrRecordNotFound" + modelName + " does not exist")
	case gorm.ErrInvalidTransaction:
		return errors.New("[gorm]:ErrInvalidTransaction" + modelName + " invalid transaction")
	default:
		return err
	}
}
