package xmatters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// -------------------------------------------------------------------------------------------------
// Group Roster Structs
// -------------------------------------------------------------------------------------------------

// GroupRoster represents a group roster in xMatters.
type GroupRoster struct {
	ID      *string         `json:"groupId"`
	Group   *GroupReference `json:"group"`
	Members []*GroupMember
}

// GroupMember represents a shorthand version of a group member.
// It contains the ID and type of the member, which can be a person, device, or group.
type GroupMember struct {
	ID         *string `json:"id" tfsdk:"id"`
	MemberType *string `json:"recipientType" tfsdk:"member_type"`
}

// GroupReference represents a shorthand version of a group in xMatters.
type GroupReference struct {
	ID            *string `json:"id,omitempty"`
	TargetName    *string `json:"targetName,omitempty"`
	RecipientType *string `json:"recipientType,omitempty"`
	GroupType     *string `json:"groupType,omitempty"`
}

// RecipientReference represents a group member in xMatters.
// Group members can be people, devices, or groups
// This object might include additional information depending on the type of group member, as defined
type RecipientReference struct {
	ID            *string `json:"id"`
	TargetName    *string `json:"targetName"`
	RecipientType *string `json:"recipientType"`

	// Optional Shared Fields
	Site        *ReferenceById    `json:"site,omitempty"`
	Supervisors *PersonPagination `json:"supervisors,omitempty"`
	Description *string           `json:"description,omitempty"`
	Observers   *RolePagination   `json:"observers,omitempty"`

	// Optional Group Fields
	AllowDuplicates        *bool              `json:"allowDuplicates,omitempty"`
	GroupType              *string            `json:"groupType,omitempty"`
	ObservedByAll          *bool              `json:"observedByAll,omitempty"`
	ResponseCount          *int64             `json:"responseCount,omitempty"`
	ResponseCountThreshold *int64             `json:"responseCountThreshold,omitempty"`
	UseDefaultDevices      *bool              `json:"useDefaultDevices,omitempty"`
	Services               *ServicePagination `json:"services,omitempty"`

	// Optional Person Fields
	FirstName   *string            `json:"firstName,omitempty"`
	Language    *string            `json:"language,omitempty"`
	LastName    *string            `json:"lastName,omitempty"`
	LicenseType *string            `json:"licenseType,omitempty"`
	PhoneLogin  *string            `json:"phoneLogin,omitempty"`
	PhonePin    *string            `json:"phonePin,omitempty"`
	Properties  *map[string]string `json:"properties,omitempty"`
	Roles       *RolePagination    `json:"roles,omitempty"`
	Timezone    *string            `json:"timezone,omitempty"`
	LastLogin   *string            `json:"lastLogin,omitempty"`
	WebLogin    *string            `json:"webLogin,omitempty"`

	// Opional Device Fields
	DefaultDevice     *bool              `json:"defaultDevice,omitempty"`
	Delay             *int64             `json:"delay,omitempty"`
	DeviceType        *string            `json:"deviceType,omitempty"`
	Name              *string            `json:"name,omitempty"`
	Owner             *PersonReference   `json:"owner,omitempty"`
	PriorityThreshold *string            `json:"priorityThreshold,omitempty"`
	Provider          *ReferenceById     `json:"provider,omitempty"`
	Sequence          *string            `json:"sequence,omitempty"`
	TestStatus        *string            `json:"testStatus,omitempty"`
	Timeframes        *[]DeviceTimeframe `json:"timeframes,omitempty"`

	// Optional DynamicTeam Fields
	ExternallyOwned *bool `json:"externallyOwned,omitempty"`
	// Criteria        *DynamicTeamCriteria `json:"criteria,omitempty"`
}

// GroupMembership represents the membership of a person, group, or device within this group.
// It contains a reference to the group and member, and for On-Call groups may optionally contain information about the specific shifts the member belongs to.
type GroupMembership struct {
	Group  GroupReference     `json:"group" validate:"required"`
	Member RecipientReference `json:"member" validate:"required"`
	Shifts ShiftPagination    `json:"shifts,omitempty"`
}

// GroupMembershipPagination contains a paginated list of group memberships.
// It extends the Pagination struct containing links to additional pages.
type GroupMembershipPagination struct {
	Pagination
	Memberships []*GroupMembership `json:"data"`
}

// -------------------------------------------------------------------------------------------------
// Group Roster Methods
// -------------------------------------------------------------------------------------------------

// GetGroupRoster retrieves the member roster of a group in xMatters.
// It requires the groupId parameter to identify the specific group, and returns a GroupRoster object.
func (xmatters *XMattersAPI) GetGroupRoster(groupId string) (GroupRoster, error) {
	uri := buildURI(fmt.Sprintf("/groups/%s/members", groupId), nil)

	// Use the GetGroupRosterPaginationSet method to get all members of the group
	groupRoster, err := xmatters.GetGroupRosterPaginationSet(uri)
	if err != nil {
		return GroupRoster{}, err
	}

	// Return the fully filled out group roster
	return groupRoster, nil
}

// GetGroupRosterPaginationSet is a recursive helper function that handles a paginated list of group rosters.
// It takes a URI as input and retrieves the paginated set from that URI.
// It checks for additional pages and recursively fetches them until all pages are retrieved.
func (xmatters *XMattersAPI) GetGroupRosterPaginationSet(uri string) (GroupRoster, error) {
	// Perform the API request.
	resp, err := xmatters.Request(http.MethodGet, uri, ContentJSON, nil)
	if err != nil {
		return GroupRoster{}, err
	}

	// Unmarshal the response body into the GroupMembershipPagination struct.
	var memberPagination GroupMembershipPagination
	err = json.Unmarshal(resp, &memberPagination)
	if err != nil {
		return GroupRoster{}, newUnmarshalError()
	}

	if len(memberPagination.Memberships) == 0 {
		return GroupRoster{}, nil
	}

	// Assign members to be returned
	var memberList []*GroupMember
	for _, member := range memberPagination.Memberships {
		memberList = append(memberList, &GroupMember{
			ID:         member.Member.ID,
			MemberType: member.Member.RecipientType,
		})
	}

	// Check for additional paginated results
	if memberPagination.Pagination.Links.Next != nil {
		nextUri := strings.ReplaceAll(*memberPagination.Pagination.Links.Next, defaultBasePath, "")
		// Use recursion to get the next set of results
		nextSet, err := xmatters.GetGroupRosterPaginationSet(nextUri)
		if err != nil {
			return GroupRoster{}, err
		}
		// Append the next set of results to the current list
		memberList = append(memberList, nextSet.Members...)
	}

	// Assign group information from the first membership entry
	groupRoster := GroupRoster{
		ID:      memberPagination.Memberships[0].Group.ID,
		Group:   &memberPagination.Memberships[0].Group,
		Members: memberList,
	}

	// Return the fully concatenated list of members from all paginated results
	return groupRoster, nil
}

// PushGroupRoster updates the members of a group in xMatters to match the desired list of members.
// This method will remove any members from the group that are not in the desired list, and add any members that are not already in the group.
// The method returns the updated group roster.
func (xmatters *XMattersAPI) PushGroupRoster(groupId string, params []*GroupMember) (GroupRoster, error) {
	currentRoster, err := xmatters.GetGroupRoster(groupId)
	if err != nil {
		return GroupRoster{}, err
	}
	// Iterate over current members and remove them from the group if they are not in the desired list
	for _, member := range currentRoster.Members {
		if !ContainsMember(*member, params) {
			if err := xmatters.DeleteGroupMembership(groupId, *member.ID); err != nil {
				return GroupRoster{}, err
			}

		}
	}
	// Iterate over desired members and add them to the group if they are not already members
	for _, member := range params {
		if !ContainsMember(*member, currentRoster.Members) {
			if _, err := xmatters.PushGroupMembership(groupId, member); err != nil {
				return GroupRoster{}, err
			}
		}
	}
	// Get the updated roster and return
	newRoster, err := xmatters.GetGroupRoster(groupId)
	if err != nil {
		return GroupRoster{}, err
	}
	return newRoster, nil
}

// DeleteGroupRoster removes all members from a group in xMatters.
// It requires the groupId parameter to identify the specific group and returns an error if any issues occur.
func (xmatters *XMattersAPI) DeleteGroupRoster(groupId string) error {
	roster, err := xmatters.GetGroupRoster(groupId)
	if err != nil {
		return err
	}
	for _, member := range roster.Members {
		if err := xmatters.DeleteGroupMembership(groupId, *member.ID); err != nil {
			return err
		}
	}

	return nil
}

// PushGroupMembership is a helper function that adds a single member to a group in xMatters.
// It requires the groupId parameter to identify the specific group and the params parameter to specify the member to be added.
// The method returns the updated GroupMember object.
// It is used internally by the PushGroupRoster method to add members to a group.
func (xmatters *XMattersAPI) PushGroupMembership(groupId string, params *GroupMember) (GroupMember, error) {
	uri := buildURI(fmt.Sprintf("/groups/%s/members", groupId), nil)

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodPost, uri, ContentJSON, params)
	if err != nil {
		return GroupMember{}, err
	}

	// Unmarshal the response body into the Recipient struct.
	var result GroupMember
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return GroupMember{}, newUnmarshalError()
	}

	// Return the paginated group members.
	return result, err
}

// DeleteGroupMembership is a helper function that removes a member from a group in xMatters.
// It requires the groupId and memberId parameters to identify the specific group and member to be removed.
// The method returns an error if any issues occur.
// It is used internally by the PushGroupRoster method to remove members from a group.
func (xmatters *XMattersAPI) DeleteGroupMembership(groupId, memberId string) error {
	uri := buildURI(fmt.Sprintf("/groups/%s/members/%s", groupId, memberId), nil) // The URI for creating or modifying a Group Member in xMatters

	// Perform the API request.
	resp, err := xmatters.Request(http.MethodDelete, uri, ContentJSON, nil)
	if err != nil {
		return err
	}

	// Unmarshal the response into a Group struct.
	var result GroupMember
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return nil
}

// ContainsMember is a helper function that checks if a GroupMember is in a given list of GroupMembers.
// It takes a GroupMember and a slice of GroupMembers as input and returns true if the member is found in the list, false otherwise.
// This function is used internally by the PushGroupRoster method to check if a member is already in the group.
func ContainsMember(member GroupMember, target []*GroupMember) bool {
	for _, m := range target {
		if *m.ID == *member.ID {
			return true
		}
	}
	return false
}
