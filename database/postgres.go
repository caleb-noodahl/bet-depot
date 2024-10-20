package database

import (
	"github.com/caleb-noodahl/bet-depot/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgErrorNotFound error
type PgErrorConflict error

type PostgresDB struct {
	Client *gorm.DB
}

func NewPostgresDB(conf *config.APIConf) (PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%v database=%s",
		conf.DBHost, conf.DBUser, conf.DBPassword, conf.DBPort, conf.DBName)

	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	return PostgresDB{
		Client: client,
	}, err

}

func (d PostgresDB) MigrateDomainModels(models ...interface{}) error {
	for _, model := range models {
		if err := d.MigrateDomainModel(&model); err != nil {
			return err
		}
	}
	return nil
}

func (d PostgresDB) MigrateDomainModel(model *interface{}) error {
	if err := d.Client.AutoMigrate(model); err != nil {
		return err
	}

	return nil
}
