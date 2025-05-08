package xmatters

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// -------------------------------------------------------------------------------------------------
// Template Structs
// -------------------------------------------------------------------------------------------------

// Template represents a template in xMatters.
type Template struct {
	StringField *string             `json:"stringField,omitempty"`
	IntField    *int64              `json:"intField,omitempty"`
	BoolField   *bool               `json:"boolField,omitempty"`
	SetField    *TemplatePagination `json:"setField,omitempty"`
	ListField   *TemplatePagination `json:"listField,omitempty"`
	ObjectField *TemplateObject     `json:"objectField,omitempty"`
}

type TemplateObject struct {
	StringField *string `json:"stringField,omitempty"`
}

// TemplatePagination represents a paginated list of templates in xMatters.
type TemplatePagination struct {
	*Pagination `json:"pagination"`
	Data        []*Template `json:"data"`
}

// -------------------------------------------------------------------------------------------------
// Method Parameter Structs
// -------------------------------------------------------------------------------------------------

// GetTemplatesParams represents parameters for getting a list of templates.
type GetTemplatesParams struct {
	Search  string `url:"search,omitempty"`
	Fields  string `url:"fields,omitempty"`
	Operand string `url:"operand,omitempty"`
	OwnedBy string `url:"ownedBy,omitempty"`
}

// PushTemplateParams represents parameters for creating or modifying a template.
type PushTemplateParams struct {
	ID         string `json:"id,omitempty"`
	TargetName string `json:"targetName,omitempty"`
}

// -------------------------------------------------------------------------------------------------
// Template Methods
// -------------------------------------------------------------------------------------------------

// GetTemplate retrieves a template in xMatters.
func (xmatters *XMattersAPI) GetTemplate(templateId *string) (Template, error) {
	uri := buildURI(fmt.Sprintf("/template/%s", *templateId), struct {
		Embed string `url:"embed,omitempty"`
	}{Embed: "templateLinks"})

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return Template{}, err
	}

	// Unmarshal the response into a Template struct.
	var result Template
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Template{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	// Return the returned Template object.
	return result, nil
}

// GetTemplateList retrieves a list of templates in xMatters.
func (xmatters *XMattersAPI) GetTemplateList(params GetTemplatesParams) ([]*Template, error) {
	uri := buildURI("/template", params) // The URI including any Query Parameters

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return []*Template{}, err
	}

	// Unmarshal the response into a TemplatePagination struct.
	var templatePag TemplatePagination
	err = json.Unmarshal(resp, &templatePag)
	if err != nil {
		return []*Template{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	// Return the TemplatePagination Data field.
	return templatePag.Data, nil
}

// PushTemplate either creates a new template or modifies an existing template in xMatters.
func (xmatters *XMattersAPI) PushTemplate(params PushTemplateParams) (Template, error) {
	uri := buildURI("/template", nil) // The URI including any Query Parameters

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodPost, uri, ContentJSON, params)
	if err != nil {
		return Template{}, err
	}

	// Unmarshal the response into a Template struct.
	var result Template
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Template{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	// Return the returned Template object.
	return result, nil
}

// DeleteTemplate deletes a template in xMatters.
func (xmatters *XMattersAPI) DeleteTemplate(templateId *string) error {
	uri := buildURI(fmt.Sprintf("/template/%s", *templateId), nil)

	// Perform the API request.
	_, err := xmatters.Request(http.MethodDelete, uri, ContentJSON, nil)
	if err != nil {
		return err
	}

	// Return
	return nil
}
