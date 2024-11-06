package dto

type CreateCategoryDTO struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	ParentID    *uint  `json:"parentId"`
	Level       int    `json:"level" binding:"required"`
	IsActive    bool   `json:"isActive"`
	SortOrder   int    `json:"sortOrder"`
}

type UpdateCategoryDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
	ParentID    *uint  `json:"parentId"`
	IsActive    *bool  `json:"isActive"`
	SortOrder   *int   `json:"sortOrder"`
}

type CategoryListDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	ParentID    *uint  `json:"parentId"`
	IsActive    bool   `json:"isActive"`
}
