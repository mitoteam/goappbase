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

	if !checkSchemaModelType(reflect.TypeOf(modelObject)) {
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

func LoadOrCreateObject[ModelT any](id any) (r *ModelT) {
	typedId, ok := mttools.AnyToInt64Ok(id)

	if !ok {
		return nil // id type is unknown
	}

	if typedId == 0 {
		return new(ModelT)
	} else {
		return LoadObject[ModelT](typedId)
	}
}

// Deletes object. Returns false if something goes wrong.
func DeleteObject(modelObject any) bool {
	var t reflect.Type

	if t, modelObject = modelObjectReflection(modelObject); t == nil {
		// not a schema model object
		return false
	}

	//make sure it is saved object (with ID set)
	v := reflect.ValueOf(modelObject).Elem() // struct itself from pointer
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

// Saves object. Returns false if something goes wrong.
func SaveObject(modelObject any) bool {
	var t reflect.Type

	if t, modelObject = modelObjectReflection(modelObject); t == nil {
		// not a schema model object
		return false
	}

	//if err := DbSchema.Db().Save(v.Interface()).Error; err != nil {
	if err := DbSchema.Db().Save(modelObject).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Query ERROR: " + err.Error())
		}

		return false
	}

	return true
}

func LoadObjectList[ModelT any]() (list []*ModelT) {
	var modelObject ModelT

	if !checkSchemaModelType(reflect.TypeOf(modelObject)) {
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

// Checks if t is type of registered model struct
func checkSchemaModelType(t reflect.Type) bool {
	if DbSchema.Db() == nil {
		//database is not opened
		return false
	}

	if !DbSchema.HasModel(t) {
		log.Printf("ERROR[checkSchemaModelType]: unknown model '%s'\n", t.String())
		return false
	}

	return true
}

// returns t  = Type of model structure, o = pointer to model object
func modelObjectReflection(modelObject any) (t reflect.Type, o any) {
	t = reflect.TypeOf(modelObject)

	//check if it is a model object pointer and dereference it's type
	if t.Kind() == reflect.Pointer {
		t = t.Elem()    // t is pointer, dereference it's type to struct
		o = modelObject // already a pointer
	} else if t.Kind() == reflect.Struct {
		o = reflect.ValueOf(modelObject).Addr().Interface() // modelObject is struct, return pointer to it
	}

	if !checkSchemaModelType(t) {
		return nil, nil
	}

	return t, o
}
