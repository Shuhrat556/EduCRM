package dto

// PaginationQuery is shared by paginated list endpoints (1-based page; omit for handler defaults).
type PaginationQuery struct {
	Page     int `form:"page" json:"page" example:"1" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" json:"page_size" example:"20" binding:"omitempty,min=1,max=100"`
}

// MessageResponse is the standard success payload for actions with no body (logout, delete).
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}
