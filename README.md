# xmatters-go

Package xmatters-go provides a Go client for interacting with the xMatters REST API.

xMatters is a communication platform that enables enterprises to manage and automate communication
with their employees, customers, and other stakeholders during incidents and other critical events.
This package allows you to programmatically access and utilize the xMatters REST API to integrate xMatters
functionality into your Go applications.

This client library is primarily used by the xMatters Terraform Provider.

[![License: MPL2.0](https://img.shields.io/badge/License-MPL2.0-yellow.svg)](./LICENSE)
[![API](https://img.shields.io/badge/API%20Docs-reference-green)](https://help.xmatters.com/xmapi/)

## Installation

```bash
go get github.com/xmatters/xmatters-go
```

## Example Usage

```go
// Create a new XMattersAPI client with your API Token
apiToken := "your-api-token"
xmattersClient, err := xmatters.NewWithAPIToken(&apiToken)
if err != nil {
    log.Fatal(err)
}

// Use the client to interact with the xMatters REST API
// For example, retrieve information about users:
users, err := xmattersClient.GetPersonList(nil)
if err != nil {
    log.Fatal(err)
}
fmt.Println(users)
```

## Available Types

### type [Device](/devices.go#L15)

`type Device struct { ... }`

Device represents a device in xMatters.

* func (*XMattersAPI) [GetDevice](/devices.go#L132)
* func (*XMattersAPI) [GetDeviceList](/devices.go#L156)
* func (*XMattersAPI) [PushDevice](/devices.go#L209)
* func (*XMattersAPI) [DeleteDevice](/devices.go#L232)

### type [Group](/groups.go#L15)

`type Group struct { ... }`

Group represents a group in xMatters.

* func (XMattersAPI) [GetGroup](/groups.go#L119)
* func (*XMattersAPI) [GetGroupList](/groups.go#L143)
* func (*XMattersAPI) [PushGroup](/groups.go#L196)
* func (*XMattersAPI) [DeleteGroup](/groups.go#L219)

### type [GroupMember](/group_roster.go#L23)

`type GroupMember struct { ... }`

GroupMember represents a shorthand version of a group member.
It contains the ID and type of the member, which can be a person, device, or group.

* func (*XMattersAPI) [GetGroupRoster](/group_roster.go#L110)
* func (*XMattersAPI) [PushGroupRoster](/group_roster.go#L179)
* func (*XMattersAPI) [PushGroupMembership](/group_roster.go#L229)
* func (*XMattersAPI) [DeleteGroupMembership](/group_roster.go#L253)

### type [Person](/people.go#L15)

`type Person struct { ... }`

Person represents a person in xMatters.

* func (*XMattersAPI) [GetPerson](/people.go#L168)
* func (*XMattersAPI) [GetPersonList](/people.go#L192)
* func (*XMattersAPI) [PushPerson](/people.go#L245)
* func (*XMattersAPI) [DeletePerson](/people.go#L268)

### type [Service](/services.go#L15)

`type Service struct { ... }`

Service represents a service in xMatters.

* func (*XMattersAPI) [GetService](/services.go#L134)
* func (*XMattersAPI) [GetServiceList](/services.go#L158)
* func (*XMattersAPI) [PushService](/services.go#L211)
* func (*XMattersAPI) [DeleteService](/services.go#L234)
* func (*XMattersAPI) [GetServiceDependency](/services.go#L253)
* func (*XMattersAPI) [PushServiceDependency](/services.go#L277)
* func (*XMattersAPI) [DeleteServiceDependency](/services.go#L300)

### type [Site](/sites.go#L15)

`type Site struct { ... }`

Site represents a site in xMatters.

* func (*XMattersAPI) [GetSite](/sites.go#L77)
* func (*XMattersAPI) [GetSiteList](/sites.go#L99)
* func (*XMattersAPI) [PushSite](/sites.go#L152)
* func (*XMattersAPI) [DeleteSite](/sites.go#L175)
