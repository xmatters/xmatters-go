package xmatters

// Pagination represents a page of results. Use the links in the links field to retrieve the rest of the result set.
type Pagination struct {
	Count *int64           `json:"count"`
	Total *int64           `json:"total"`
	Links *PaginationLinks `json:"links"`
}

// PaginationLinks provides links to the current, previous, and next pages of a paginated result set.
type PaginationLinks struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Self     *string `json:"self"`
}

// ReferenceById represents the identifier of a resource.
type ReferenceById struct {
	ID *string `json:"id"`
}

// ReferenceByName identifies a resource by name.
type ReferenceByName struct {
	Name *string `json:"name"`
}
