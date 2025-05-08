package xmatters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// -------------------------------------------------------------------------------------------------
// Site Structs
// -------------------------------------------------------------------------------------------------

// Site represents a site in xMatters.
type Site struct {
	Address1   *string  `json:"address1,omitempty"`
	Address2   *string  `json:"address2,omitempty"`
	City       *string  `json:"city,omitempty"`
	Country    *string  `json:"country,omitempty"`
	ID         *string  `json:"id"`
	Language   *string  `json:"language,omitempty"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
	Name       *string  `json:"name,omitempty"`
	PostalCode *string  `json:"postalCode,omitempty"`
	State      *string  `json:"state,omitempty"`
	Status     *string  `json:"status,omitempty"`
	Timezone   *string  `json:"timezone,omitempty"`
}

// SitePagination contains a paginated list of sites.
// It extends the Pagination struct containing links to additional pages.
type SitePagination struct {
	*Pagination
	Sites []*Site `json:"data,omitempty"`
}

// -------------------------------------------------------------------------------------------------
// Method Parameter Structs
// -------------------------------------------------------------------------------------------------

// GetSitesParams contains available API query parameters for the GetSiteList method.
type GetSitesParams struct {
	Search   string `url:"search,omitempty"`
	Operand  string `url:"operand,omitempty"`
	Fields   string `url:"fields,omitempty"`
	Country  string `url:"country,omitempty"`
	Geocoded *bool  `url:"geocoded,omitempty"`
	Status   string `url:"status,omitempty"`
}

// GetSiteParams contains available API body parameters for the PushSite method.
type PushSiteParams struct {
	// Required Fields
	Name     string `json:"name"`
	Country  string `json:"country"`
	Language string `json:"language"`
	Timezone string `json:"timezone"`
	// Optional Fields
	Address1   *string  `json:"address1"`
	Address2   *string  `json:"address2"`
	City       *string  `json:"city"`
	ID         string   `json:"id,omitempty"`
	Latitude   *float64 `json:"latitude"`
	Longitude  *float64 `json:"longitude"`
	PostalCode *string  `json:"postalCode"`
	State      *string  `json:"state"`
	Status     string   `json:"status,omitempty"`
}

// -------------------------------------------------------------------------------------------------
// Site Methods
// -------------------------------------------------------------------------------------------------

// GetSite retrieves a site in xMatters.
// It requires the siteId parameter to identify the specific site, and returns a Site object.
func (xmatters *XMattersAPI) GetSite(siteId string) (Site, error) {
	uri := buildURI(fmt.Sprintf("/sites/%s", siteId), nil)

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return Site{}, err
	}

	// Unmarshal the response into a Site struct.
	var result Site
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Site{}, newUnmarshalError()
	}

	// Return the returned Site object.
	return result, nil
}

// GetSiteList retrieves a list of sites in xMatters.
// It accepts optional query parameters to filter the results and returns a slice of Site objects.
func (xmatters *XMattersAPI) GetSiteList(params GetSitesParams) ([]*Site, error) {
	uri := buildURI("/sites", params) // The URI including any Query Parameters

	// Use the GetSitePaginationSet method to get all paginated results
	siteList, err := xmatters.GetSitePaginationSet(uri)
	if err != nil {
		return []*Site{}, err
	}

	// Return the full list of Sites.
	return siteList, nil
}

// GetSitePaginationSet is a recursive helper function that handles a paginated list of sites.
// It takes a URI as input and retrieves the paginated set from that URI.
// It checks for additional pages and recursively fetches them until all pages are retrieved.
func (xmatters *XMattersAPI) GetSitePaginationSet(uri string) ([]*Site, error) {
	// Perform the API request with provided URI
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return []*Site{}, err
	}

	// Unmarshal the response into a SitePagination struct.
	var sitePagination SitePagination
	err = json.Unmarshal(resp, &sitePagination)
	if err != nil {
		return []*Site{}, newUnmarshalError()
	}

	// Assign first page of sites to be returned
	siteList := sitePagination.Sites

	// Check for additional paginated results
	if sitePagination.Pagination.Links.Next != nil {
		// Remove defaultBasePath (/api/xm/1) from the next URI
		nextUri := strings.ReplaceAll(*sitePagination.Pagination.Links.Next, defaultBasePath, "")
		// Use recursion to get the next set of results
		nextSet, err := xmatters.GetSitePaginationSet(nextUri)
		if err != nil {
			return []*Site{}, err
		}
		siteList = append(siteList, nextSet...)
	}

	// Return the fully concatenated list of sites from all paginated results
	return siteList, nil
}

// PushSite either creates a new site in xMatters or modifies an existing site.
// It requires the PushSiteParams struct containing the site details.
// It returns the created or modified Site object.
// If the params.ID is provided it updates the existing site; otherwise, it creates a new one.
func (xmatters *XMattersAPI) PushSite(params PushSiteParams) (Site, error) {
	uri := buildURI("/sites", nil) // The URI including any Query Parameters

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodPost, uri, ContentJSON, params)
	if err != nil {
		return Site{}, err
	}

	// Unmarshal the response into a Site struct.
	var result Site
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Site{}, newUnmarshalError()
	}

	// Return the returned Site object.
	return result, nil
}

// DeleteSite deletes a site in xMatters.
// It requires the siteId parameter to identify the specific site to be deleted.
// It returns an error if the deletion fails.
func (xmatters *XMattersAPI) DeleteSite(siteId *string) error {
	uri := buildURI(fmt.Sprintf("/sites/%s", *siteId), nil)

	// Perform the API request.
	_, err := xmatters.Request(http.MethodDelete, uri, ContentJSON, nil)
	if err != nil {
		return err
	}

	// Return
	return nil
}
