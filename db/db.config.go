package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseClient struct {
	MongoClient *mongo.Client
	PostgresDB  *gorm.DB
}

var dbClient *DatabaseClient

func InitDBClient(appConfig AppConfig) error {
	if appConfig.Dbtype == "mongo" {
		clientOptions := options.Client().ApplyURI(appConfig.DbUri).
			SetMaxPoolSize(uint64(appConfig.DbMaxPoolSize)).
			SetMinPoolSize(uint64(appConfig.DbMinPoolSize)).
			SetMaxConnIdleTime(10 * time.Minute).
			SetConnectTimeout(10 * time.Second)

		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			return fmt.Errorf("failed to connect to mongodb: %v", err)
		}

		// Retry Logic
		for i := 0; i < 3; i++ {
			if err = client.Ping(context.Background(), nil); err == nil {
				log.Println("Successfully connected to MongoDB")
				dbClient = &DatabaseClient{MongoClient: client}
				return nil
			}
			time.Sleep(2 * time.Second)
		}
		return fmt.Errorf("MongoDB connection failed after retries: %v", err)
	} else if appConfig.Dbtype == "postgres" {
		dbUri := fmt.Sprintf("%s/%s?sslmode=disable", appConfig.DbUri, appConfig.DbName)

		// Initialize GORM connection
		db, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent), // Adjust logging level as needed
		})
		if err != nil {
			return fmt.Errorf("failed to connect to postgres: %v", err)
		}

		// Get SQL DB connection for connection pooling settings
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get SQL DB from GORM: %v", err)
		}

		sqlDB.SetMaxOpenConns(appConfig.DbMaxPoolSize)
		sqlDB.SetMaxIdleConns(appConfig.DbMinPoolSize)
		sqlDB.SetConnMaxIdleTime(10 * time.Minute)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)

		// Ping to check connection
		if err = sqlDB.Ping(); err != nil {
			return fmt.Errorf("failed to ping postgres: %v", err)
		}

		log.Println("Successfully connected to PostgreSQL and schema migrated")
		dbClient = &DatabaseClient{PostgresDB: db}
		return nil
	}
	return fmt.Errorf("unsupported database type: %s", appConfig.Dbtype)
}

func GetDBClient() *DatabaseClient {
	return dbClient
}

func GetMongoDataBase(dbName string) *mongo.Database {
	return dbClient.MongoClient.Database(dbName)
}

func GetPostgresDB() *gorm.DB {
	return dbClient.PostgresDB
}

func AutoMigrateModels(models ...interface{}) error {
	db := GetPostgresDB()
	if db == nil {
		return fmt.Errorf("PostgreSQL database is not initialized")
	}

	err := db.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("auto migration failed: %v", err)
	}
	log.Println("Auto migration successful for provided models")
	return nil
}
