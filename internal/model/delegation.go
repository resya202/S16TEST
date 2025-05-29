package model

import "time"

// HourlyDelegation maps to the delegation_hourly table
type HourlyDelegation struct {
	ID            uint      `gorm:"primaryKey;column:id"`
	ValidatorAddr string    `gorm:"column:validator_addr"`
	DelegatorAddr string    `gorm:"column:delegator_addr"`
	Timestamp     time.Time `gorm:"column:timestamp"`
	AmountUAtom   int64     `gorm:"column:amount_uatom"`
	ChangeUAtom   int64     `gorm:"column:change_uatom"`
}

func (HourlyDelegation) TableName() string {
	return "delegation_hourly"
}
