// Package model defines data struct based on DB schema
package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type Todo struct {
	ID          string    `gorm:"type:varchar(255);not null"`
	UserID      string    `gorm:"type:varchar(255);not null"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Finished    bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"type:timestamp"`
	UpdatedAt   time.Time `gorm:"type:timestamp"`
}

type DBClient interface {
	StartTransaction() (DBClient, error)
	GetConnection() *gorm.DB
	Rollback() error
	Commit() error

	GetTodoList(filter map[string]interface{}, offset int, limit int) ([]Todo, error)
	GetTodo(id string) (Todo, error)
	DeleteTodo(id string) error
	CreateTodo(t Todo) (string, error)
	UpdateTodo(t Todo) error
}

type DB struct {
	dbConnect *gorm.DB
}

var db *DB

// DBConn creates a mysql connection.
// Connection string 'conStr' should be 'user:pwd@tcp(host:port)/dbname'.
func DBInit(conStr string) {
	if db == nil {
		var err error
		dbConnection, err := gorm.Open("mysql", conStr+"?charset=utf8&parseTime=True&loc=Local")
		// db connection will be closed if there's no request for a while
		// which would cause 500 error when a new request comes.
		// diable pool here to avoid this problem.
		dbConnection.DB().SetMaxIdleConns(0)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Faile to create db connection pool")
		} else {
			log.WithFields(log.Fields{
				"message": dbConnection.GetErrors(),
				"db":      conStr,
			}).Info("connected to mysql")
		}
		db = &DB{dbConnection}
	}
	db.dbConnect.SetLogger(log.StandardLogger())
	// db.Debug message will be logged be logrug with Info level
	db.dbConnect.Debug().AutoMigrate(&Todo{})
}

func New() *DB {
	return db
}

func (d *DB) GetConnection() *gorm.DB {
	return d.dbConnect
}

func (d *DB) StartTransaction() (DBClient, error) {
	tx := d.dbConnect.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &DB{tx}, nil
}

func (d *DB) Rollback() error {
	return d.dbConnect.Rollback().Error
}

func (d *DB) Commit() error {
	return d.dbConnect.Commit().Error
}

func (d *DB) GetTodoList(filter map[string]interface{}, offset int, limit int) (l []Todo, e error) {
	dbConnect := d.GetConnection()
	e = dbConnect.Where(filter).Offset(offset).
		Limit(limit).
		Find(&l).Error
	return
}

func (d *DB) GetTodo(id string) (t Todo, e error) {
	e = d.dbConnect.Where("id = ?", id).First(&t).Error
	return
}

func (d *DB) DeleteTodo(id string) (e error) {
	var t Todo
	e = d.dbConnect.Where("id = ?", id).First(&t).Error
	if e != nil {
		return
	}
	e = d.dbConnect.Debug().Delete(&t).Error
	return
}

func (d *DB) CreateTodo(t Todo) (string, error) {
	t.ID = uuid.New().String()
	t.CreatedAt = time.Now().UTC()
	t.UpdatedAt = t.CreatedAt
	dbConnect := d.GetConnection()
	err := dbConnect.Debug().Create(&t).Error
	return t.ID, err
}

func (d *DB) UpdateTodo(t Todo) error {
	t.UpdatedAt = time.Now().UTC()
	dbConnect := d.GetConnection()
	err := dbConnect.Debug().Save(&t).Error
	return err
}
