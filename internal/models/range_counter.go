package models

type RangeCounter struct {
	ID       uint  `gorm:"primaryKey;column:id"`
	LastUsed int64 `gorm:"not null;column:last_used"`
}

func (RangeCounter) TableName() string {
	return "range_counters"
}
