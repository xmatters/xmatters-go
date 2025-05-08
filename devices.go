package xmatters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// -------------------------------------------------------------------------------------------------
// Device Structs
// -------------------------------------------------------------------------------------------------

// Device represents a device in xMatters.
type Device struct {
	ID                *string            `json:"id"`
	TargetName        *string            `json:"targetName,omitempty"`
	Country           *string            `json:"country,omitempty"`
	DefaultDevice     *bool              `json:"defaultDevice,omitempty"`
	Delay             *int32             `json:"delay,omitempty"`
	DeviceType        *string            `json:"deviceType"`
	EmailAddress      *string            `json:"emailAddress,omitempty"`
	ExternalKey       *string            `json:"externalKey,omitempty"`
	ExternallyOwned   *bool              `json:"externallyOwned,omitempty"`
	Name              *string            `json:"name"`
	Owner             *PersonReference   `json:"owner"`
	PhoneNumber       *string            `json:"phoneNumber,omitempty"`
	PIN               *string            `json:"pin,omitempty"`
	PriorityThreshold *string            `json:"priorityThreshold,omitempty"`
	Sequence          *int32             `json:"sequence,omitempty"`
	Status            *string            `json:"status,omitempty"`
	TestStatus        *string            `json:"testStatus,omitempty"`
	Timeframes        []*DeviceTimeframe `json:"timeframes,omitempty"`
	TwoWayDevice      *bool              `json:"twoWayDevice,omitempty"`
}

// DevicePagination contains a paginated list of devices.
// It extends the Pagination struct containing links to additional pages.
type DevicePagination struct {
	Devices []*Device `json:"data,omitempty"`
	*Pagination
}

// -----------------------------------------------------------------------------------
// DeviceTimeframe Structs
// -----------------------------------------------------------------------------------

// DeviceTimeframe represents a timeframe during which a device is active and able to receive notifications.
type DeviceTimeframe struct {
	Name              *string   `json:"name" tfsdk:"name"`
	StartTime         *string   `json:"startTime" tfsdk:"start_time"`
	DurationInMinutes *int32    `json:"durationInMinutes" tfsdk:"duration_in_minutes"`
	Days              []*string `json:"days" tfsdk:"days"`
	ExcludeHolidays   *bool     `json:"excludeHolidays" tfsdk:"exclude_holidays"`
}

// DeviceTimeframePagination contains a paginated list of device timeframes.
// It extends the Pagination struct containing links to additional pages.
type DeviceTimeframePagination struct {
	*Pagination
	Data []*DeviceTimeframe `json:"data,omitempty"`
}

// -------------------------------------------------------------------------------------------------
// Method Parameter Structs
// -------------------------------------------------------------------------------------------------

// GetDevicesParams contains available API query parameters for the GetDeviceList method.
type GetDevicesParams struct {
	Embed        string `url:"embed,omitempty"`
	DeviceStatus string `url:"deviceStatus,omitempty"`
	DeviceType   string `url:"deviceType,omitempty"`
	DeviceNames  string `url:"deviceNames,omitempty"`
}

// PushDeviceParams contains available API body parameters for the PushDevice method.
type PushDeviceParams struct {
	// Required Fields
	DeviceType        string             `json:"deviceType"`
	Name              string             `json:"name"`
	Owner             string             `json:"owner"`
	Sequence          *int32             `json:"sequence"`
	PriorityThreshold string             `json:"priorityThreshold"`
	TestStatus        string             `json:"testStatus"`
	Timeframes        []*DeviceTimeframe `json:"timeframes"`
	// Optional Fields
	ID              string  `json:"id,omitempty"`
	Country         string  `json:"country,omitempty"`
	DefaultDevice   *bool   `json:"defaultDevice,omitempty"`
	Delay           *int32  `json:"delay"`
	EmailAddress    string  `json:"emailAddress,omitempty"`
	ExternalKey     *string `json:"externalKey"`
	ExternallyOwned *bool   `json:"externallyOwned"`
	PhoneNumber     string  `json:"phoneNumber,omitempty"`
	PIN             string  `json:"pin,omitempty"`
	Status          string  `json:"status,omitempty"`
	TwoWayDevice    *bool   `json:"twoWayDevice"`
}

// -------------------------------------------------------------------------------------------------
// Device Methods
// -------------------------------------------------------------------------------------------------

// Custom Unmarshaller for Device to handle embedded timeframes
// This is necessary because the JSON structure for timeframes is nested within a pagination object.
func (d *Device) UnmarshalJSON(data []byte) error {
	// Define an alias to avoid recursion
	type Alias Device
	aux := &struct {
		Timeframes struct {
			Data []*DeviceTimeframe `json:"data"`
		} `json:"timeframes"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	// Unmarshal the JSON into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal Device: %w", err)
	}

	// Assign the extracted device timeframes
	d.Timeframes = aux.Timeframes.Data

	return nil
}

// GetDevice retrieves a device in xMatters.
// It requires the deviceId parameter to identify the specific device, and returns a Device object.
// A URL parameter is added to the request URI to embed timeframes of the device in the response.
func (xmatters *XMattersAPI) GetDevice(deviceId string) (Device, error) {
	uri := buildURI(fmt.Sprintf("/devices/%s", deviceId), struct {
		Embed string `url:"embed"`
	}{Embed: "timeframes"})

	// Perform the API request
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return Device{}, err
	}

	// Unmarshal the response into a Device struct.
	var result Device
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Device{}, newUnmarshalError()
	}

	// Return the details of the specific Device.
	return result, nil
}

// GetDeviceList retrieves a list of devices in xMatters.
// It accepts optional query parameters to filter the results and returns a slice of Device objects.
func (xmatters *XMattersAPI) GetDeviceList(params GetDevicesParams) ([]*Device, error) {
	uri := buildURI("/devices", params) // The URI including the given Query Parameters

	// Use the GetDevicePaginationSet method to get all paginated results
	deviceList, err := xmatters.GetDevicePaginationSet(uri)
	if err != nil {
		return []*Device{}, err
	}

	// Return the list of devices
	return deviceList, nil
}

// GetDevicePaginationSet is a recursive helper function that handles a paginated list of devices.
// It takes a URI as input and retrieves the paginated set from that URI.
// It checks for additional pages and recursively fetches them until all pages are retrieved.
func (xmatters *XMattersAPI) GetDevicePaginationSet(uri string) ([]*Device, error) {
	// Perform the API request with provided URI
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return []*Device{}, err
	}

	// Unmarshal the response into a DevicePagination struct.
	var devicePagination DevicePagination
	err = json.Unmarshal(resp, &devicePagination)
	if err != nil {
		return []*Device{}, newUnmarshalError()
	}

	// Assign devices to be returned
	deviceList := devicePagination.Devices

	// Check for additional paginated results
	if devicePagination.Pagination.Links.Next != nil {
		// Remove defaultBasePath (/api/xm/1) from the next URI
		nextUri := strings.ReplaceAll(*devicePagination.Pagination.Links.Next, defaultBasePath, "")
		// Use recursion to get the next set of results
		nextSet, err := xmatters.GetDevicePaginationSet(nextUri)
		if err != nil {
			return []*Device{}, err
		}
		deviceList = append(deviceList, nextSet...)
	}

	// Return the fully concatenated list of devices from all paginated results
	return deviceList, nil
}

// PushDevice either creates a new device in xMatters or modifies an existing device.
// It requires the PushDeviceParams struct containing the device details.
// It returns the created or modified Device object.
// If the params.ID is provided it updates the existing device; otherwise, it creates a new one.
func (xmatters *XMattersAPI) PushDevice(params PushDeviceParams) (Device, error) {
	uri := buildURI("/devices", nil) // The URI for creating or modifying a Device in xMatters

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodPost, uri, ContentJSON, params)
	if err != nil {
		return Device{}, err
	}

	// Unmarshal the response into a Device struct.
	var result Device
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return Device{}, newUnmarshalError()
	}

	// Return the created or Modified Device details.
	return result, nil
}

// DeleteDevice deletes a device in xMatters.
// It requires the deviceId parameter to identify the specific device to be deleted.
// It returns an error if the deletion fails.
func (xmatters *XMattersAPI) DeleteDevice(params string) error {
	uri := buildURI(fmt.Sprintf("/devices/%s", params), nil) // The URI for Deleting a Device in xMatters

	resp, err := xmatters.Request(http.MethodDelete, uri, ContentJSON, nil)
	if err != nil {
		return err
	}

	// Unmarshal the response into a Device struct.
	var result Device
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	// Return the deleted Device details.
	return nil
}
