package xmatters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// -------------------------------------------------------------------------------------------------
// Service Structs
// -------------------------------------------------------------------------------------------------

// Service represents a service in xMatters.
type Service struct {
	ID              *string         `json:"id"`
	TargetName      *string         `json:"targetName,omitempty"`
	RecipientType   *string         `json:"recipientType,omitempty"`
	ServiceType     *string         `json:"serviceType,omitempty"`
	ServiceTier     *string         `json:"serviceTier,omitempty"`
	Description     *string         `json:"description,omitempty"`
	ServiceLinks    []*ServiceLink  `json:"serviceLinks"`
	OwnedBy         *GroupReference `json:"ownedBy,omitempty"`
	ExternallyOwned *bool           `json:"externallyOwned,omitempty"`
	Status          *string         `json:"status,omitempty"`
}

// ServicePagination contains a paginated list of services.
// It extends the Pagination struct containing links to additional pages.
type ServicePagination struct {
	*Pagination
	Services []*Service `json:"data,omitempty"`
}

// ServiceLink represents an optional URL link with associated label assigned to an xMatters service.
type ServiceLink struct {
	Label *string `json:"label" tfsdk:"link_text"`
	URL   *string `json:"url" tfsdk:"url"`
}

// ServiceLinksPagination contains a paginated list of service links.
// It extends the Pagination struct containing links to additional pages.
type ServiceLinksPagination struct {
	*Pagination
	Data []*ServiceLink `json:"data,omitempty"`
}

// -------------------------------------------------------------------------------------------------
// Service Dependancy Structs
// -------------------------------------------------------------------------------------------------

// ServiceDependency represents a service dependency relationship in xMatters.
type ServiceDependency struct {
	ID               *string           `json:"id"`
	Service          *ServiceReference `json:"service"`
	DependentService *ServiceReference `json:"dependentService"`
}

// ServiceDependencyPagination contains a paginated list of service dependencies.
// It extends the Pagination struct containing links to additional pages.
type ServiceDependencyPagination struct {
	*Pagination
	Data []*ServiceDependency `json:"data"`
}

// ServiceReference represents a shorthand version of a service in xMatters.
type ServiceReference struct {
	ID         *string `json:"id"`
	TargetName *string `json:"targetName,omitempty"`
}

// -------------------------------------------------------------------------------------------------
// Method Parameter Structs
// -------------------------------------------------------------------------------------------------

// GetServicesParams contains available API query parameters for the GetServiceList method.
type GetServicesParams struct {
	Search  string `url:"search,omitempty"`
	Fields  string `url:"fields,omitempty"`
	Operand string `url:"operand,omitempty"`
	OwnedBy string `url:"ownedBy,omitempty"`
}

// PushServiceParams contains available API body parameters for the PushService method.
type PushServiceParams struct {
	ID           string          `json:"id,omitempty"`
	TargetName   string          `json:"targetName"`
	Description  *string         `json:"description"`
	ServiceType  string          `json:"serviceType"`
	ServiceTier  *string         `json:"serviceTier"`
	OwnedBy      *GroupReference `json:"ownedBy"`
	ServiceLinks []*ServiceLink  `json:"serviceLinks"`
}

// PushServiceDependencyParams contains available API body parameters for the PushServiceDependency method.
type PushServiceDependencyParams struct {
	ID                 string `json:"id"`
	ServiceID          string `json:"serviceId"`
	DependentServiceID string `json:"dependentServiceId"`
}

// -------------------------------------------------------------------------------------------------
// Service Methods
// -------------------------------------------------------------------------------------------------

// Custom Unmarshaller for Service to handle embedded service links
// This is necessary because the JSON structure for service links is nested within a pagination object.
func (s *Service) UnmarshalJSON(data []byte) error {
	// Define an alias to avoid recursion
	type Alias Service
	aux := &struct {
		ServiceLinks struct {
			Data []*ServiceLink `json:"data"`
		} `json:"serviceLinks"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	// Unmarshal the JSON into the auxiliary struct
	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("failed to unmarshal Service: %w", err)
	}

	// Assign the extracted service links
	s.ServiceLinks = aux.ServiceLinks.Data

	return nil
}

// GetService retrieves a service in xMatters.
// It requires the serviceId parameter to identify the specific service, and returns a Service object.
// A URL parameter is added to the request URI to embed service links of the service in the response.
func (xmatters *XMattersAPI) GetService(serviceId string) (Service, error) {
	uri := buildURI(fmt.Sprintf("/services/%s", serviceId), struct {
		Embed string `url:"embed"`
	}{Embed: "serviceLinks"})

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return Service{}, err
	}

	// Unmarshal the response into a Service struct.
	var result Service
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Service{}, newUnmarshalError()
	}

	// Return the returned Service object.
	return result, nil
}

// GetServiceList retrieves a list of services in xMatters.
// It accepts optional query parameters to filter the results and returns a slice of Service objects.
func (xmatters *XMattersAPI) GetServiceList(params GetServicesParams) ([]*Service, error) {
	uri := buildURI("/services", params) // The URI including any Query Parameters

	// Use the GetServicePaginationSet method to get all paginated results
	serviceList, err := xmatters.GetServicePaginationSet(uri)
	if err != nil {
		return []*Service{}, err
	}

	// Return the full list of Services.
	return serviceList, nil
}

// GetServicePaginationSet is a recursive helper function that handles a paginated list of services.
// It takes a URI as input and retrieves the paginated set from that URI.
// It checks for additional pages and recursively fetches them until all pages are retrieved.
func (xmatters *XMattersAPI) GetServicePaginationSet(uri string) ([]*Service, error) {
	// Perform the API request with provided URI
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return []*Service{}, err
	}

	// Unmarshal the response into a ServicePagination struct.
	var servicePagination ServicePagination
	err = json.Unmarshal(resp, &servicePagination)
	if err != nil {
		return []*Service{}, newUnmarshalError()
	}

	// Assign services to be returned
	serviceList := servicePagination.Services

	// Check for additional paginated results
	if servicePagination.Pagination.Links.Next != nil {
		// Remove defaultBasePath (/api/xm/1) from the next URI
		nextUri := strings.ReplaceAll(*servicePagination.Pagination.Links.Next, defaultBasePath, "")
		// Use recursion to get the next set of results
		nextSet, err := xmatters.GetServicePaginationSet(nextUri)
		if err != nil {
			return []*Service{}, err
		}
		serviceList = append(serviceList, nextSet...)
	}

	// Return the fully concatenated list of services from all paginated results
	return serviceList, nil
}

// PushService either creates a new service in xMatters or modifies an existing service.
// It requires the PushServiceParams struct containing the service details.
// It returns the created or modified Service object.
// If the params.ID is provided it updates the existing service; otherwise, it creates a new one.
func (xmatters *XMattersAPI) PushService(params PushServiceParams) (Service, error) {
	uri := buildURI("/services", nil) // The URI including any Query Parameters

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodPost, uri, ContentJSON, params)
	if err != nil {
		return Service{}, err
	}

	// Unmarshal the response into a Service struct.
	var result Service
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Service{}, newUnmarshalError()
	}

	// Return the returned Service object.
	return result, nil
}

// DeleteService deletes a service in xMatters.
// It requires the serviceId parameter to identify the specific service to be deleted.
// It returns an error if the deletion fails.
func (xmatters *XMattersAPI) DeleteService(serviceId string) error {
	uri := buildURI(fmt.Sprintf("/services/%s", serviceId), nil)

	// Perform the API request.
	_, err := xmatters.Request(http.MethodDelete, uri, ContentJSON, nil)
	if err != nil {
		return err
	}

	// Return
	return nil
}

// -------------------------------------------------------------------------------------------------
// Service Dependancy Methods
// -------------------------------------------------------------------------------------------------

// GetServiceDependency retrieves a service dependency in xMatters.
// It requires the dependencyId parameter to identify the specific service dependency, and returns a ServiceDependency object.
func (xmatters *XMattersAPI) GetServiceDependency(dependencyId string) (ServiceDependency, error) {
	uri := buildURI(fmt.Sprintf("/service-dependencies/%s", dependencyId), nil) // The URI including any Query Parameters

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return ServiceDependency{}, err
	}

	// Unmarshal the response into a ServiceDependency struct.
	var result ServiceDependency
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return ServiceDependency{}, newUnmarshalError()
	}

	// Return the returned ServiceDependency object.
	return result, err
}

// PushServiceDependency either creates a new service dependency in xMatters or modifies an existing service dependency.
// It requires the PushServiceDependencyParams struct containing the service dependency details.
// It returns the created or modified ServiceDependency object.
// If the params.ID is provided it updates the existing service dependency; otherwise, it creates a new one.
func (xmatters *XMattersAPI) PushServiceDependency(params PushServiceDependencyParams) (ServiceDependency, error) {
	uri := buildURI("/service-dependencies", nil) // The URI including any Query Parameters

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodPost, uri, ContentJSON, params)
	if err != nil {
		return ServiceDependency{}, err
	}

	// Unmarshal the response into a ServiceDependency struct.
	var result ServiceDependency
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return ServiceDependency{}, newUnmarshalError()
	}

	// Return the returned ServiceDependency object.
	return result, err
}

// DeleteServiceDependency deletes a service dependency in xMatters.
// It requires the serviceDepId parameter to identify the specific service dependency to be deleted.
// It returns an error if the deletion fails.
func (xmatters *XMattersAPI) DeleteServiceDependency(serviceDepId string) error {
	uri := buildURI(fmt.Sprintf("/service-dependencies/%s", serviceDepId), nil) // The URI including any Query Parameters

	// Perform the API request.
	_, err := xmatters.Request(http.MethodDelete, uri, ContentJSON, nil)
	if err != nil {
		return err
	}

	// Return
	return nil
}
