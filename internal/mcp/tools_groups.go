// Package mcp provides the MCP server implementation for ABAP ADT tools.
// tools_groups.go defines tool group definitions for selective disablement.
package mcp

// toolGroups defines groups of tools that can be selectively disabled.
// Short codes: 5/U=UI5, T=Tests, C=CTS
func toolGroups() map[string][]string {
	groups := map[string][]string{
		"5": { // UI5/BSP tools (also mapped as "U")
			"UI5ListApps", "UI5GetApp", "UI5GetFileContent",
		},
		"T": { // Test tools
			"RunUnitTests", "RunATCCheck",
		},
		"C": { // CTS/Transport tools
			"ListTransports", "GetTransport",
			"CreateTransport", "ReleaseTransport", "DeleteTransport",
		},
	}
	// Map "U" to same tools as "5"
	groups["U"] = groups["5"]
	return groups
}
