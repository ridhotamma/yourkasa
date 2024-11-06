package dto

import "time"

type CreateGroupDTO struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	ImageURL    string     `json:"imageUrl"`
	IsActive    bool       `json:"isActive"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	GroupType   string     `json:"groupType" binding:"required"`
	SortOrder   int        `json:"sortOrder"`
}

type UpdateGroupDTO struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ImageURL    string     `json:"imageUrl"`
	IsActive    *bool      `json:"isActive"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	GroupType   string     `json:"groupType"`
	SortOrder   *int       `json:"sortOrder"`
}
