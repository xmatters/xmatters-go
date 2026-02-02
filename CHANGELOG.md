# Changelog

## [0.3.0] - TBD

FEATURES:

- Added SearchCriteria and SearchCriterion types for dynamic group filtering.
- Added Criteria field to Group struct for dynamic group membership configuration.
- Added Criteria field to PushGroupParams for creating and updating dynamic groups.
- Renamed group_roster.go to group_members.go for improved clarity and consistency.

ENHANCEMENTS:

- Enhanced GetGroupsParams with additional filter parameters:
  - CreatedAfter: Filter groups created after a specific date.
  - CreatedBefore: Filter groups created before a specific date.
  - CreatedFrom: Filter groups created from a specific date.
  - CreatedTo: Filter groups created to a specific date.
  - MemberLicenseType: Filter groups by member license type.
- Added custom UnmarshalJSON method for SearchCriteria to handle nested criterion pagination.
- Updated PushGroupParams Description field to be required (non-omitempty).

CHANGES:

- Renamed GroupRoster related types and methods to GroupMembers for better API alignment.
- Restructured group membership code organization.

## [0.2.0] - 2025-09-09

FEATURES:

- Added Device support with full CRUD operations.
- Added DeviceTimeframe type for device scheduling.

ENHANCEMENTS:

- Added comprehensive device-related API methods.
- Added device filtering and pagination support.

## [0.1.0] - 2025-05-22

- Initial release.
- Support for Person, Group, Service, Site, and Role resources.
- Basic CRUD operations for supported resources.
- Authentication via API Token or Basic Auth.
