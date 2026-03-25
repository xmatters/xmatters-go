package xmatters

import (
	"encoding/json"
	"testing"
)

func TestSearchCriteriaUnmarshalJSON_PaginatedCriterion(t *testing.T) {
	payload := []byte(`{
		"operand": "AND",
		"criterion": {
			"count": 1,
			"data": [
				{
					"criterionType": "BASIC_FIELD",
					"field": "NAME",
					"operand": "CONTAINS",
					"value": "March"
				}
			]
		}
	}`)

	var criteria SearchCriteria
	if err := json.Unmarshal(payload, &criteria); err != nil {
		t.Fatalf("expected no error unmarshalling paginated criterion, got: %v", err)
	}

	if criteria.Operand == nil || *criteria.Operand != "AND" {
		t.Fatalf("expected operand AND, got: %#v", criteria.Operand)
	}
	if len(criteria.Criterion) != 1 {
		t.Fatalf("expected 1 criterion, got %d", len(criteria.Criterion))
	}
	if criteria.Criterion[0].Field == nil || *criteria.Criterion[0].Field != "NAME" {
		t.Fatalf("expected criterion field NAME, got: %#v", criteria.Criterion[0].Field)
	}
}

func TestSearchCriteriaUnmarshalJSON_ArrayCriterion(t *testing.T) {
	payload := []byte(`{
		"operand": "OR",
		"criterion": [
			{
				"criterionType": "BASIC_FIELD",
				"field": "NAME",
				"operand": "CONTAINS",
				"value": "March"
			}
		]
	}`)

	var criteria SearchCriteria
	if err := json.Unmarshal(payload, &criteria); err != nil {
		t.Fatalf("expected no error unmarshalling criterion array, got: %v", err)
	}

	if criteria.Operand == nil || *criteria.Operand != "OR" {
		t.Fatalf("expected operand OR, got: %#v", criteria.Operand)
	}
	if len(criteria.Criterion) != 1 {
		t.Fatalf("expected 1 criterion, got %d", len(criteria.Criterion))
	}
}

func TestSearchCriteriaUnmarshalJSON_NullCriterion(t *testing.T) {
	payload := []byte(`{
		"operand": "AND",
		"criterion": null
	}`)

	var criteria SearchCriteria
	if err := json.Unmarshal(payload, &criteria); err != nil {
		t.Fatalf("expected no error unmarshalling null criterion, got: %v", err)
	}

	if criteria.Operand == nil || *criteria.Operand != "AND" {
		t.Fatalf("expected operand AND, got: %#v", criteria.Operand)
	}
	if criteria.Criterion != nil {
		t.Fatalf("expected nil criterion for null payload, got: %#v", criteria.Criterion)
	}
}

func TestSearchCriteriaUnmarshalJSON_InvalidCriterionType(t *testing.T) {
	payload := []byte(`{
		"operand": "AND",
		"criterion": "unexpected"
	}`)

	var criteria SearchCriteria
	if err := json.Unmarshal(payload, &criteria); err == nil {
		t.Fatalf("expected error unmarshalling invalid criterion payload")
	}
}
