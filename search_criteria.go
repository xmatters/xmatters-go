package xmatters

import (
	"bytes"
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
	Data []*SearchCriterion `json:"data,omitempty"`
}

// -------------------------------------------------------------------------------------------------
// Search Criteria Methods
// -------------------------------------------------------------------------------------------------

// Custom Unmarshaller for SearchCriteria to handle criterion payload variants.
// Read responses provide criterion in a pagination object ({"data": [...]}) while
// other payloads may provide criterion as a direct array. This implementation accepts both.
func (c *SearchCriteria) UnmarshalJSON(data []byte) error {
	var aux struct {
		Operand   *string         `json:"operand,omitempty"`
		Criterion json.RawMessage `json:"criterion,omitempty"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal SearchCriteria: %w", err)
	}

	c.Operand = aux.Operand
	c.Criterion = nil

	rawCriterion := bytes.TrimSpace(aux.Criterion)
	if len(rawCriterion) == 0 || bytes.Equal(rawCriterion, []byte("null")) {
		return nil
	}

	switch rawCriterion[0] {
	case '{':
		var criterionPagination SearchCriterionPagination
		if err := json.Unmarshal(rawCriterion, &criterionPagination); err != nil {
			return fmt.Errorf("failed to unmarshal SearchCriteria criterion pagination: %w", err)
		}
		c.Criterion = criterionPagination.Data
	case '[':
		var criterion []*SearchCriterion
		if err := json.Unmarshal(rawCriterion, &criterion); err != nil {
			return fmt.Errorf("failed to unmarshal SearchCriteria criterion array: %w", err)
		}
		c.Criterion = criterion
	default:
		return fmt.Errorf("failed to unmarshal SearchCriteria criterion: unexpected JSON token %q", rawCriterion[0])
	}

	return nil
}
