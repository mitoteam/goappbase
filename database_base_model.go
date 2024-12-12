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
	typedId, ok := mttools.AnyToInt64Ok(id)

	if !ok || typedId == 0 { //id is empty
		return nil
	}

	var modelObject ModelT

	if !checkSchemaModel(reflect.TypeOf(modelObject)) {
		return nil
	}

	if err := DbSchema.Db().First(&modelObject, typedId).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Query ERROR: " + err.Error())
		}

		return nil
	}

	return &modelObject
}

// Deletes object. Returns false if something goes wrong.
func DeleteObject(modelObject any) bool {
	//check if it is a model object pointer and unpack it
	t := reflect.TypeOf(modelObject)
	v := reflect.ValueOf(modelObject)

	if t.Kind() == reflect.Pointer {
		v = v.Elem()
		t = t.Elem()
		modelObject = v.Interface()
	}

	if !checkSchemaModel(t) {
		return false
	}

	//make sure there is saved object (with ID set)
	if v.FieldByName(reflect.TypeFor[BaseModel]().Name()).FieldByName("ID").Int() == 0 {
		return false
	}

	if err := DbSchema.Db().Delete(modelObject).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Query ERROR: " + err.Error())
		}

		return false
	}

	return true
}

func LoadObjectList[ModelT any]() (list []*ModelT) {
	var modelObject ModelT

	if !checkSchemaModel(reflect.TypeOf(modelObject)) {
		return []*ModelT{} //empty list
	}

	if err := DbSchema.Db().Model(&modelObject).Find(&list).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Query ERROR: " + err.Error())
		}

		return []*ModelT{} //empty list
	}

	return list
}

func checkSchemaModel(t reflect.Type) bool {
	if DbSchema.Db() == nil {
		//database is not opened
		return false
	}

	if !DbSchema.HasModel(t) {
		log.Printf("ERROR[checkSchemaModel]: unknown model '%s'\n", t.String())
		return false
	}

	return true
}

// func UniqueSlice[S ~[]E, E any](slice S) S {
