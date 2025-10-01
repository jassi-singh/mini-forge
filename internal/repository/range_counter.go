package repository

import (
	"github.com/jassi-singh/mini-forge/internal/models"
	"github.com/jassi-singh/mini-forge/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RangeCounterRepository interface {
	GetAndIncrement() (int64, error)
}

type rangeCounterRepo struct {
	db     *gorm.DB
	config *utils.Config
}

func NewRangeCounterRepository(db *gorm.DB, config *utils.Config) RangeCounterRepository {
	return &rangeCounterRepo{db: db, config: config}
}

func (r *rangeCounterRepo) GetAndIncrement() (int64, error) {
	var counter int64

	txErr := r.db.Transaction(func(tx *gorm.DB) error {
		var rangeCounter models.RangeCounter
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&rangeCounter, 1).Error; err != nil {
			return err
		}

		counter = rangeCounter.LastUsed

		rangeCounter.LastUsed += int64(r.config.RangeSize)

		if err := tx.Save(&rangeCounter).Error; err != nil {
			return err
		}

		return nil
	})

	return counter, txErr
}
