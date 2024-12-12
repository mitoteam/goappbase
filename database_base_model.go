package goappbase

import (
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/mitoteam/mttools"
	gorm "gorm.io/gorm"
)

// gorm.Model alternative without DeletedAt column (to disable Soft Delete)
// see https://gorm.io/docs/delete.html#Soft-Delete
type BaseModel struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func LoadObject[ModelT any](id any) (r *ModelT) {
	if DbSchema.Db() == nil {
		//database is not opened
		return nil
	}

	typedId, ok := mttools.AnyToInt64Ok(id)

	if !ok && typedId == 0 {
		//id is empty
		return nil
	}

	var modelObject ModelT
	t := reflect.TypeOf(modelObject)

	if !DbSchema.HasModel(t) {
		log.Printf("LoadObject ERROR: unknown model '%s'\n", t.Name())
		return nil
	}

	if err := DbSchema.Db().First(&modelObject, typedId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			log.Println("Query ERROR: " + err.Error())
			return nil
		}
	}

	return &modelObject
}

// func UniqueSlice[S ~[]E, E any](slice S) S {
