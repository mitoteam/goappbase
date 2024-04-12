package goappbase

import (
	"log"

	"github.com/glebarez/sqlite"
	gorm "gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const dbFileName = "data.db"

type DbSchemaType struct {
	schema []interface{}
	db     *gorm.DB
}

var DbSchema *DbSchemaType

func init() {
	DbSchema = &DbSchemaType{}

	DbSchema.schema = make([]interface{}, 0)
}

func (schema *DbSchemaType) AddModel(model interface{}) {
	schema.schema = append(schema.schema, model)
}

func (schema *DbSchemaType) GetDb() *gorm.DB {
	return schema.db
}

func (db_schema *DbSchemaType) Open() (*gorm.DB, error) {
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
	for i := 0; i < len(db_schema.schema); i++ {
		db_schema.db.AutoMigrate(db_schema.schema[i])
	}

	log.Println("Database migration done")

	return db_schema.db, nil
}

func (schema *DbSchemaType) Close() {
	sqlDB, err := schema.db.DB()

	if err != nil {
		sqlDB.Close()
	}

	schema.db = nil
}
