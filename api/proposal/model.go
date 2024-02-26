package proposal

import (
	"time"

	"gorm.io/gorm"
)

type Proposal struct {
	// taken from from gorm.Model, so we can json strigify properly
	ID        uint64         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Document string `json:"document"`
	Start    uint32 `json:"start"`
	End      uint32 `json:"end"`
}
