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
	db gorm.DB
}

func NewRangeCounterRepository(db gorm.DB) RangeCounterRepository {
	return &rangeCounterRepo{db: db}
}

func (r *rangeCounterRepo) GetAndIncrement() (int64, error) {
	var counter int64

	txErr := r.db.Transaction(func(tx *gorm.DB) error {
		var rangeCounter models.RangeCounter
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&rangeCounter, 1).Error; err != nil {
			return err
		}

		counter = rangeCounter.LastUsed + 1
		rangeCounter.LastUsed += int64(utils.RANGE_SIZE)

		if err := tx.Save(&rangeCounter).Error; err != nil {
			return err
		}

		return nil
	})

	return counter, txErr
}
