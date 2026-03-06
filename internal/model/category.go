package model

import (
	"time"
)

type Category struct {
	ID          string        `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string        `gorm:"type:varchar(50);uniqueIndex"`
	Description *string       `gorm:"type:text"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	AuctionItem []AuctionItem `gorm:"foreignKey:CategoryID;references:ID"`
}

type CategoryCreateRequest struct {
	Name        string `form:"name" validate:"required"`
	Description string `form:"name,omitempty" validate:"omitempty"`
}

type CategoryUpdateRequest struct {
	Name        string `form:"name,omitempty" validate:"omitempty"`
	Description string `form:"name,omitempty" validate:"omitempty"`
}

type CategoryResource struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
}

func (c *Category) Resource() CategoryResource {
	return CategoryResource{
		ID:          c.ID,
		Name:        c.Name,
		Description: *c.Description,
	}
}

func CategoryResources(categories []Category) []CategoryResource {
	if len(categories) == 0 {
		return []CategoryResource{}
	}

	responses := make([]CategoryResource, 0, len(categories))
	for i := range categories {
		responses = append(responses, categories[i].Resource())
	}
	return responses
}
