package xmatters

import (
	"strings"
	"testing"
)

func TestBuildURI_PreservesDelimiterEncodingForSpecialQueryParams(t *testing.T) {
	tests := []struct {
		name      string
		uri       string
		contains  string
		notEquals string
	}{
		{
			name: "groups",
			uri: buildURI("/people", GetPeopleParams{
				Groups: "Group+One,Group+Two",
			}),
			contains:  "groups=Group+One,Group+Two",
			notEquals: "groups=Group+One%2CGroup+Two",
		},
		{
			name: "deviceNames",
			uri: buildURI("/devices", GetDevicesParams{
				DeviceNames: "Phone+One,Phone+Two",
			}),
			contains:  "deviceNames=Phone+One,Phone+Two",
			notEquals: "deviceNames=Phone+One%2CPhone+Two",
		},
		{
			name: "sites",
			uri: buildURI("/groups", GetGroupsParams{
				Sites: "Default+Site,Great+White+North",
			}),
			contains:  "sites=Default+Site,Great+White+North",
			notEquals: "sites=Default+Site%2CGreat+White+North",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(tt.uri, tt.contains) {
				t.Fatalf("expected query param serialization in %q, got %q", tt.contains, tt.uri)
			}

			if strings.Contains(tt.uri, tt.notEquals) {
				t.Fatalf("expected comma delimiter to remain unescaped for %s, got %q", tt.name, tt.uri)
			}
		})
	}
}
