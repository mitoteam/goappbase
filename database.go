package goappbase

import (
	"log"
	"reflect"

	"github.com/glebarez/sqlite"
	gorm "gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const dbFileName = "data.db"

type dbSchemaType struct {
	modelMap map[string]any // name = typename, value = empty struct of this type
	db       *gorm.DB
}

var DbSchema *dbSchemaType

func init() {
	DbSchema = &dbSchemaType{}

	DbSchema.modelMap = make(map[string]any, 0)
}

func (schema *dbSchemaType) AddModel(modelType reflect.Type) {
	//ensure it is a struct
	if modelType.Kind() != reflect.Struct {
		log.Panicf("modelType %s is not a struct", modelType.Name())
	}

	schema.modelMap[modelType.String()] = reflect.New(modelType).Elem().Interface()
}

func (schema *dbSchemaType) HasModel(modelType reflect.Type) bool {
	_, exists := schema.modelMap[modelType.String()]
	return exists
}

func (schema *dbSchemaType) Db() *gorm.DB {
	return schema.db
}

func (db_schema *dbSchemaType) Open() (*gorm.DB, error) {
	var err error

	db_schema.db, err = gorm.Open(sqlite.Open(dbFileName), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})

	if err != nil {
		return nil, err
	}

	log.Printf("Database %s opened\n", dbFileName)

	// Migrate the schema
	//log.Printf("DBG: %+v\n", db_schema.modelMap)
	for _, modelObject := range db_schema.modelMap {
		db_schema.db.AutoMigrate(modelObject)
	}

	log.Printf("Database migration done (schema model count: %d)\n", len(db_schema.modelMap))

	return db_schema.db, nil
}

func (schema *dbSchemaType) Close() {
	sqlDB, err := schema.db.DB()

	if err != nil {
		sqlDB.Close()
	}

	log.Printf("Database %s closed\n", dbFileName)

	schema.db = nil
}
