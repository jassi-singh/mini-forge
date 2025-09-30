package models

type RangeCounter struct {
	ID       uint `gorm:"primaryKey;column:id"`
	LastUsed uint `gorm:"not null;column:last_used"`
}

func (RangeCounter) TableName() string {
	return "range_counters"
}
