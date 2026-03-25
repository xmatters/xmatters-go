package xmatters

import (
	"encoding/json"
	"fmt"
)

// -------------------------------------------------------------------------------------------------
// Search Criteria Structs
// -------------------------------------------------------------------------------------------------

// SearchCriteria represents the search criteria for filtering groups in xMatters.
type SearchCriteria struct {
	Operand   *string            `json:"operand,omitempty" tfsdk:"operand"`
	Criterion []*SearchCriterion `json:"criterion,omitempty"`
}

// SearchCriterion represents a single search criterion used in filtering groups.
type SearchCriterion struct {
	CriterionType *string `json:"criterionType,omitempty" tfsdk:"criterion_type"`
	Field         *string `json:"field,omitempty" tfsdk:"field"`
	Operand       *string `json:"operand,omitempty" tfsdk:"operand"`
	Value         *string `json:"value,omitempty" tfsdk:"value"`
}

// SearchCriterionPagination contains a paginated list of search criterion.
type SearchCriterionPagination struct {
	*Pagination
	SearchCriterion []*SearchCriterion `json:"data,omitempty"`
}

// -------------------------------------------------------------------------------------------------
// Search Criteria Methods
// -------------------------------------------------------------------------------------------------

// Custom Unmarshaller for SearchCriteria to handle embedded criterion array.
// This is necessary because the JSON structure for these fields are nested within pagination objects.
func (c *SearchCriteria) UnmarshalJSON(data []byte) error {
	type Alias SearchCriteria
	aux := &struct {
		Criterion struct {
			Criterion []*SearchCriterion `json:"data"`
		} `json:"criterion,omitempty" tfsdk:"criterion"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	// Unmarshal the JSON into the auxiliary struct
	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("failed to unmarshal Search Criterion: %w", err)
	}

	// Assign the extracted attributes
	c.Criterion = aux.Criterion.Criterion
	return nil
}
