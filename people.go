package xmatters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// -------------------------------------------------------------------------------------------------
// People Structs
// -------------------------------------------------------------------------------------------------

// Person represents a person in xMatters.
type Person struct {
	ID              *string        `json:"id"`
	TargetName      *string        `json:"targetName"`
	FirstName       *string        `json:"firstName"`
	LastName        *string        `json:"lastName"`
	Roles           []*Role        `json:"roles"`
	Status          *string        `json:"status,omitempty"`
	WebLogin        *string        `json:"webLogin,omitempty"`
	Site            *ReferenceById `json:"site,omitempty"`
	Timezone        *string        `json:"timezone,omitempty"`
	Language        *string        `json:"language,omitempty"`
	Supervisors     []*Person      `json:"supervisors,omitempty"`
	PhoneLogin      *string        `json:"phoneLogin,omitempty"`
	LicenseType     *string        `json:"licenseType,omitempty"`
	ExternalKey     *string        `json:"externalKey,omitempty"`
	ExternallyOwned *bool          `json:"externallyOwned,omitempty"`
	LastLogin       *string        `json:"lastLogin,omitempty"`
}

// PersonPagination contains a paginated list of people.
// It extends the Pagination struct containing links to additional pages.
type PersonPagination struct {
	*Pagination
	People []*Person `json:"data"`
}

// PersonReference represents a shorthand version of a person in xMatters.
type PersonReference struct {
	ID         *string `json:"id"`
	TargetName *string `json:"targetName"`
	FirstName  *string `json:"firstName"`
	LastName   *string `json:"lastName"`
}

// -------------------------------------------------------------------------------------------------
// User Quota Structs
// -------------------------------------------------------------------------------------------------

// UserQuotas represents the user licence quotas for an xMatters instance.
type UserQuotas struct {
	StakeholderUsersEnabled *bool         `json:"stakeholderUsersEnabled"`
	StakeholderUsers        *QuotaDetails `json:"stakeholderUsers"`
	FullUsers               *QuotaDetails `json:"fullUsers"`
}

// QuotaDetails represents the details of the quotas applied to User Type for an xMatters instance.
type QuotaDetails struct {
	Total  *int64 `json:"total"`
	Active *int64 `json:"active"`
	Unused *int64 `json:"unused"`
}

// -------------------------------------------------------------------------------------------------
// Method Parameter Structs
// -------------------------------------------------------------------------------------------------

// GetPeopleParams contains available API query parameters for the GetPersonList method.
type GetPeopleParams struct {
	Embed string `url:"embed,omitempty"`
	// Provider Search Object
	Terms   string `url:"search,omitempty"`
	Fields  string `url:"fields,omitempty"`
	Operand string `url:"operand,omitempty"`
	// Provider Filters Object
	CreatedAfter       string `url:"createdAfter,omitempty"`
	CreatedBefore      string `url:"createdBefore,omitempty"`
	CreatedFrom        string `url:"createdFrom,omitempty"`
	CreatedTo          string `url:"createdTo,omitempty"`
	DevicesExists      *bool  `url:"devices.exists,omitempty"`
	DevicesEmailExists *bool  `url:"devices.email.exists,omitempty"`
	DevicesFailsafe    *bool  `url:"devices.failsafe.exists,omitempty"`
	DevicesMobile      *bool  `url:"devices.mobile.exists,omitempty"`
	DevicesSMS         *bool  `url:"devices.sms.exists,omitempty"`
	DevicesVoice       *bool  `url:"devices.voice.exists,omitempty"`
	DevicesStatus      string `url:"devices.status,omitempty"`
	DevicesTestStatus  string `url:"devices.testStatus,omitempty"`
	EmailAddress       string `url:"emailAddress,omitempty"`
	FirstName          string `url:"firstName,omitempty"`
	Groups             string `url:"groups,omitempty"`
	GroupsExists       *bool  `url:"groups.exists,omitempty"`
	LastName           string `url:"lastName,omitempty"`
	LicenseType        string `url:"licenseType,omitempty"`
	PhoneNumber        string `url:"phoneNumber,omitempty"`
	Roles              string `url:"roles,omitempty"`
	Site               string `url:"site,omitempty"`
	Status             string `url:"status,omitempty"`
	Supervisors        string `url:"supervisors,omitempty"`
	SupervisorsExists  *bool  `url:"supervisors.exists,omitempty"`
	TargetName         string `url:"targetName,omitempty"`
	WebLogin           string `url:"webLogin,omitempty"`
	// Provider Options Object
	SortBy    string `url:"sortBy,omitempty"`
	SortOrder string `url:"sortOrder,omitempty"`
}

// PushPersonParams contains available API body parameters for the PushPerson method.
type PushPersonParams struct {
	// Required Fields
	TargetName  string    `json:"targetName"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Roles       []*string `json:"roles"`
	LicenseType string    `json:"licenseType"`
	Site        string    `json:"site"`
	Language    string    `json:"language"`
	Supervisors []*string `json:"supervisors"`
	Timezone    string    `json:"timezone"`
	WebLogin    string    `json:"webLogin"`
	// Optional Fields
	ID              string  `json:"id,omitempty"`
	Status          string  `json:"status,omitempty"`
	PhoneLogin      *string `json:"phoneLogin"`
	PhonePin        string  `json:"phonePin,omitempty"`
	ExternalKey     *string `json:"externalKey"`
	ExternallyOwned *bool   `json:"externallyOwned"`
}

// -------------------------------------------------------------------------------------------------
// People Methods
// -------------------------------------------------------------------------------------------------

// Custom Unmarshaller for Person to handle embedded roles and supervisors
// This is necessary because the JSON structure for roles and supervisors is nested within pagination objects.
func (p *Person) UnmarshalJSON(data []byte) error {
	// Define an alias to avoid recursion
	type Alias Person
	aux := &struct {
		Roles struct {
			Roles []*Role `json:"data"`
		} `json:"roles"`
		Supervisors struct {
			People []*Person `json:"data"`
		} `json:"supervisors,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	// Unmarshal the JSON into the auxiliary struct
	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("failed to unmarshal Person: %w", err)
	}

	// Assign the extracted attributes
	p.Roles = aux.Roles.Roles
	p.Supervisors = aux.Supervisors.People

	return nil
}

// GetPerson retrieves a person in xMatters.
// It requires the personId parameter to identify the specific person, and returns a Person object.
// A URL parameter is added to the request URI to embed the roles and supervisors of the person in the response.
func (xmatters *XMattersAPI) GetPerson(personId string) (Person, error) {
	uri := buildURI(fmt.Sprintf("/people/%s", personId), struct {
		Embed string `url:"embed"`
	}{Embed: "roles,supervisors"})

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return Person{}, err
	}

	// Unmarshal the response into an Person struct.
	var result Person
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Person{}, newUnmarshalError()
	}

	// Return the returned Person object.
	return result, nil
}

// GetPersonList retrieves a list of people in xMatters.
// It accepts optional query parameters to filter the results and returns a slice of Person objects.
func (xmatters *XMattersAPI) GetPersonList(params GetPeopleParams) ([]*Person, error) {
	uri := buildURI("/people", params)

	// Use the GetPersonPaginationSet method to get all paginated results
	personList, err := xmatters.GetPersonPaginationSet(uri)
	if err != nil {
		return []*Person{}, err
	}

	// Return the full list of People.
	return personList, nil
}

// GetPersonPaginationSet is a recursive helper function that handles a paginated list of people.
// It takes a URI as input and retrieves the paginated set from that URI.
// It checks for additional pages and recursively fetches them until all pages are retrieved.
func (xmatters *XMattersAPI) GetPersonPaginationSet(uri string) ([]*Person, error) {
	// Perform the API request with provided URI
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return []*Person{}, err
	}

	// Unmarshal the response into a PersonPagination struct.
	var personPagination PersonPagination
	err = json.Unmarshal(resp, &personPagination)
	if err != nil {
		return []*Person{}, newUnmarshalError()
	}

	// Assign people to be returned
	personList := personPagination.People

	// Check for additional paginated results
	if personPagination.Pagination.Links.Next != nil {
		// Remove defaultBasePath (/api/xm/1) from the next URI
		nextUri := strings.ReplaceAll(*personPagination.Pagination.Links.Next, defaultBasePath, "")
		// Use recursion to get the next set of results
		nextSet, err := xmatters.GetPersonPaginationSet(nextUri)
		if err != nil {
			return []*Person{}, err
		}
		personList = append(personList, nextSet...)
	}

	// Return the fully concatenated list of people from all paginated results
	return personList, nil
}

// PushPerson either creates a new person in xMatters or modifies an existing person.
// It requires the PushPersonParams struct containing the person details.
// It returns the created or modified Person object.
// If the params.ID is provided it updates the existing person; otherwise, it creates a new one.
func (xmatters *XMattersAPI) PushPerson(params PushPersonParams) (Person, error) {
	uri := buildURI("/people", nil) // The URI for creating or modifying a Person in xMatters

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodPost, uri, ContentJSON, params)
	if err != nil {
		return Person{}, err
	}

	// Unmarshal the response into a Person struct.
	var result Person
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Person{}, newUnmarshalError()
	}

	// Return the created or Modified Device details.
	return result, nil
}

// DeletePerson deletes a person in xMatters.
// It requires the personId parameter to identify the specific person to be deleted.
// It returns an error if the deletion fails.
func (xmatters *XMattersAPI) DeletePerson(personId *string) error {
	uri := buildURI(fmt.Sprintf("/people/%s", *personId), nil)

	// Perform the API request.
	_, err := xmatters.Request(http.MethodDelete, uri, ContentJSON, nil)
	if err != nil {
		return err
	}

	// Return
	return nil
}

// -------------------------------------------------------------------------------------------------
// User Quota Methods
// -------------------------------------------------------------------------------------------------

// GetUserQuotas retrieves the user license quotas for an xMatters instance.
func (xmatters *XMattersAPI) GetUserQuotas() (UserQuotas, error) {
	uri := buildURI("/people/license-quotas", nil)

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return UserQuotas{}, err
	}

	// Unmarshal the response into an UserQuotas struct.
	var result UserQuotas
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return UserQuotas{}, newUnmarshalError()
	}

	// Return the returned UserQuotas object.
	return result, nil
}
