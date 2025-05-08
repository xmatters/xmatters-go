package xmatters

// Role represents a role in xMatters.
type Role struct {
	ID          *string `json:"id,omitempty"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// RolePagination contains a paginated list of roles.
// It extends the Pagination struct containing links to additional pages.
type RolePagination struct {
	*Pagination
	Roles []*Role `json:"data,omitempty"`
}
