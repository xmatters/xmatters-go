package xmatters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// -------------------------------------------------------------------------------------------------
// Group Structs
// -------------------------------------------------------------------------------------------------

// Group represents a group in xMatters.
type Group struct {
	ID                *string            `json:"id"`
	TargetName        *string            `json:"targetName"`
	Status            *string            `json:"status"`
	Description       *string            `json:"description,omitempty"`
	GroupType         *string            `json:"groupType,omitempty"`
	AllowDuplicates   *bool              `json:"allowDuplicates,omitempty"`
	Timezone          *string            `json:"timezone,omitempty"`
	Site              *ReferenceById     `json:"site,omitempty"`
	ObservedByAll     *bool              `json:"observedByAll,omitempty"`
	Observers         []*ReferenceByName `json:"observers,omitempty"`
	UseDefaultDevices *bool              `json:"useDefaultDevices,omitempty"`
	Supervisors       []*ReferenceById   `json:"supervisors,omitempty"`
	Services          []*Service         `json:"services,omitempty"`
	ExternalKey       *string            `json:"externalKey,omitempty"`
	ExternallyOwned   *bool              `json:"externallyOwned,omitempty"`
}

// GroupPagination contains a paginated list of groups.
// It extends the Pagination struct containing links to additional pages.
type GroupPagination struct {
	*Pagination
	Groups []*Group `json:"data,omitempty"`
}

// -------------------------------------------------------------------------------------------------
// Method Parameter Structs
// -------------------------------------------------------------------------------------------------

// GetGroupsParams contains available API query parameters for the GetGroupList method.
type GetGroupsParams struct {
	Embed string `url:"embed,omitempty"`
	// Provider Search Object
	Terms   string `url:"search,omitempty"`
	Fields  string `url:"fields,omitempty"`
	Operand string `url:"operand,omitempty"`
	// Provider Filters Object
	GroupType    string `url:"groupType,omitempty"`
	MemberExists string `url:"member.exists,omitempty"`
	Members      string `url:"members,omitempty"`
	Sites        string `url:"sites,omitempty"`
	Status       string `url:"status,omitempty"`
	Supervisors  string `url:"supervisors,omitempty"`
	// Provider Options Object
	SortBy    string `url:"sortBy,omitempty"`
	SortOrder string `url:"sortOrder,omitempty"`
}

// PushGroupParams contains available API body parameters for the PushGroup method.
type PushGroupParams struct {
	ID                string             `json:"id,omitempty"`
	TargetName        string             `json:"targetName"`
	AllowDuplicates   *bool              `json:"allowDuplicates,omitempty"`
	Description       string             `json:"description,omitempty"`
	ExternalKey       string             `json:"externalKey,omitempty"`
	ExternallyOwned   *bool              `json:"externallyOwned,omitempty"`
	GroupType         string             `json:"groupType,omitempty"`
	ObservedByAll     *bool              `json:"observedByAll,omitempty"`
	Observers         []*ReferenceByName `json:"observers,omitempty"`
	Site              string             `json:"site,omitempty"`
	Status            string             `json:"status,omitempty"`
	UseDefaultDevices *bool              `json:"useDefaultDevices,omitempty"`
	Supervisors       []*ReferenceById   `json:"supervisors,omitempty"`
}

// -------------------------------------------------------------------------------------------------
// Group Methods
// -------------------------------------------------------------------------------------------------

// Custom Unmarshaller for Group to handle embedded observers, supervisors, and services
// This is necessary because the JSON structure for these fields are nested within pagination objects.
func (g *Group) UnmarshalJSON(data []byte) error {
	type Alias Group
	aux := &struct {
		Observers struct {
			Observers []*ReferenceByName `json:"data"`
		} `json:"observers"`
		Supervisors struct {
			Supervisors []*ReferenceById `json:"data"`
		} `json:"supervisors"`
		Services struct {
			Services []*Service `json:"data"`
		} `json:"services"`
		*Alias
	}{
		Alias: (*Alias)(g),
	}

	// Unmarshal the JSON into the auxiliary struct
	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("failed to unmarshal Group: %w", err)
	}

	// Assign the extracted attributes
	g.Observers = aux.Observers.Observers
	g.Supervisors = aux.Supervisors.Supervisors
	g.Services = aux.Services.Services

	return nil
}

// GetGroup retrieves a group in xMatters.
// It requires the groupId parameter to identify the specific group, and returns a Group object.
// A URL parameter is added to the request URI to embed the supervisors, observers, and services.
func (xmatters XMattersAPI) GetGroup(groupId string) (Group, error) {
	uri := buildURI(fmt.Sprintf("/groups/%s", groupId), struct {
		Embed string `url:"embed"`
	}{Embed: "supervisors,observers,services"})

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return Group{}, err
	}

	// Unmarshal the response into an Group struct.
	var result Group
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Group{}, newUnmarshalError()
	}

	// Return the returned Group object.
	return result, nil
}

// GetGroupList retrieves a list of groups in xMatters.
// It accepts optional query parameters to filter the results and returns a slice of Group objects.
func (xmatters *XMattersAPI) GetGroupList(params GetGroupsParams) ([]*Group, error) {
	uri := buildURI("/groups", params)

	// Use the GetGroupPaginationSet method to retrieve all paginated results
	groupList, err := xmatters.GetGroupPaginationSet(uri)
	if err != nil {
		return []*Group{}, err
	}

	// Return the full list of Groups.
	return groupList, nil
}

// GetGroupPaginationSet is a recursive helper function that handles a paginated list of groups.
// It takes a URI as input and retrieves the paginated set from that URI.
// It checks for additional pages and recursively fetches them until all pages are retrieved.
func (xmatters *XMattersAPI) GetGroupPaginationSet(uri string) ([]*Group, error) {
	// Perform the API request with provided URI
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return []*Group{}, err
	}

	// Unmarshal the response into a GroupPagination struct.
	var groupPagination GroupPagination
	err = json.Unmarshal(resp, &groupPagination)
	if err != nil {
		return []*Group{}, newUnmarshalError()
	}

	// Assign groups to be returned
	groupList := groupPagination.Groups

	// Check for additional paginated results
	if groupPagination.Pagination.Links.Next != nil {
		// Remove defaultBasePath (/api/xm/1) from the next URI
		nextUri := strings.ReplaceAll(*groupPagination.Pagination.Links.Next, defaultBasePath, "")
		// Use recursion to get the next set of results
		nextSet, err := xmatters.GetGroupPaginationSet(nextUri)
		if err != nil {
			return []*Group{}, err
		}
		groupList = append(groupList, nextSet...)
	}

	// Return the fully concatenated list of groups from all paginated results
	return groupList, nil
}

// PushGroup either creates a new group in xMatters or modifies an existing group.
// It requires the PushGroupParams struct to specify the group details.
// It returns the created or modified Group object.
// If the params.ID is provided it updates the existing group; otherwise, it creates a new one.
func (xmatters *XMattersAPI) PushGroup(params PushGroupParams) (Group, error) {
	uri := buildURI("/groups", nil) // The URI for creating or modifying a Group in xMatters

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodPost, uri, ContentJSON, params)
	if err != nil {
		return Group{}, err
	}

	// Unmarshal the response into a Group struct.
	var result Group
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Group{}, newUnmarshalError()
	}

	// Return the created or Modified Device details.
	return result, nil
}

// DeleteGroup deletes a group in xMatters.
// It requires the groupId parameter to identify the specific group to be deleted.
// It returns an error if the deletion fails.
func (xmatters *XMattersAPI) DeleteGroup(groupId string) error {
	uri := buildURI(fmt.Sprintf("/groups/%s", groupId), nil) // The URI for Deleting a Group in xMatters

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodDelete, uri, ContentJSON, nil)
	if err != nil {
		return err
	}

	// Unmarshal the response into a Group struct.
	var result Group
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	// Return the deleted Group details.
	return nil
}
