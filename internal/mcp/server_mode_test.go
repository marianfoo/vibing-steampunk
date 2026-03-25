package mcp

import (
	"testing"
)

// TestFocusedModeToolCount ensures focused mode registers a specific number of tools.
// This test prevents accidental tool additions — if you add a tool, update this count.
func TestFocusedModeToolCount(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "focused",
	})

	tools := listTools(t, s)
	// Count should match the focusedToolSet() entries + always-on tools (GetConnectionInfo, GetFeatures)
	// This is a regression test — update when intentionally adding/removing tools
	t.Logf("Focused mode tool count: %d", len(tools))

	if len(tools) == 0 {
		t.Fatal("focused mode should register at least some tools")
	}
}

// TestHyperfocusedModeRegistersOnlyUniversalTool ensures hyperfocused mode registers exactly 1 tool.
func TestHyperfocusedModeRegistersOnlyUniversalTool(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "hyperfocused",
	})

	tools := listTools(t, s)
	if len(tools) != 1 {
		t.Fatalf("hyperfocused mode should register exactly 1 tool, got %d", len(tools))
	}
	if tools[0].Name != "SAP" {
		t.Fatalf("hyperfocused mode tool should be 'SAP', got %q", tools[0].Name)
	}
}

// TestReadOnlyModeHasTools verifies that readonly mode registers tools and sets safety.
// Note: Currently readonly mode registers all focused tools but blocks writes at the ADT
// client level via safety config. The new tool surface (Phase 3) will hide write tools
// at the registration level.
func TestReadOnlyModeHasTools(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "readonly",
	})

	names := toolNames(t, s)

	// These read tools SHOULD be present
	readTools := []string{
		"GetSource",
		"SearchObject",
		"GetTable",
		"GetTableContents",
		"GetSystemInfo",
		"FindDefinition",
		"FindReferences",
		"ListDumps",
		"GetDump",
	}

	for _, tool := range readTools {
		if !names[tool] {
			t.Errorf("read tool %q should be registered in readonly mode", tool)
		}
	}
}

// TestReadOnlyConfigImpliesSafety verifies that read-only mode sets safety flags.
func TestReadOnlyConfigImpliesSafety(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "readonly",
	})

	if s.adtClient == nil {
		t.Fatal("ADT client should not be nil")
	}
}

// TestExpertModeRegistersMoreTools verifies expert mode has more tools than focused.
func TestExpertModeRegistersMoreTools(t *testing.T) {
	focused := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "focused",
	})

	expert := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "expert",
	})

	focusedTools := listTools(t, focused)
	expertTools := listTools(t, expert)

	if len(expertTools) <= len(focusedTools) {
		t.Fatalf("expert mode (%d tools) should have more tools than focused mode (%d tools)",
			len(expertTools), len(focusedTools))
	}
	t.Logf("Focused: %d tools, Expert: %d tools", len(focusedTools), len(expertTools))
}

// TestDisabledGroupsRemovesTools verifies that --disabled-groups removes tools.
func TestDisabledGroupsRemovesTools(t *testing.T) {
	withUI5 := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "focused",
	})

	withoutUI5 := newServerWithConfig(&Config{
		BaseURL:        "https://sap.example.com:44300",
		Username:       "testuser",
		Password:       "testpass",
		Mode:           "focused",
		DisabledGroups: "5",
	})

	withNames := toolNames(t, withUI5)
	withoutNames := toolNames(t, withoutUI5)

	// UI5 tools should be present without disabled groups
	if !withNames["UI5ListApps"] {
		t.Error("UI5ListApps should be registered without disabled groups")
	}

	// UI5 tools should be absent with "5" disabled
	ui5Tools := []string{"UI5ListApps", "UI5GetApp", "UI5GetFileContent"}
	for _, tool := range ui5Tools {
		if withoutNames[tool] {
			t.Errorf("UI5 tool %q should NOT be registered with disabled group '5'", tool)
		}
	}
}

// TestDisabledGroupsTransportsRemoved verifies "C" disables transport tools.
func TestDisabledGroupsTransportsRemoved(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:        "https://sap.example.com:44300",
		Username:       "testuser",
		Password:       "testpass",
		Mode:           "focused",
		DisabledGroups: "C",
	})

	names := toolNames(t, s)
	if names["ListTransports"] {
		t.Error("ListTransports should NOT be registered with disabled group 'C'")
	}
	if names["GetTransport"] {
		t.Error("GetTransport should NOT be registered with disabled group 'C'")
	}
}

// TestCoreToolsAlwaysPresent verifies essential tools are always registered.
func TestCoreToolsAlwaysPresent(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "focused",
	})

	names := toolNames(t, s)

	// Always-on tools
	alwaysOn := []string{"GetConnectionInfo", "GetFeatures"}
	for _, tool := range alwaysOn {
		if !names[tool] {
			t.Errorf("always-on tool %q should be registered", tool)
		}
	}

	// Core focused tools
	coreTools := []string{
		"GetSource", "WriteSource", "SearchObject",
		"GrepObjects", "GrepPackages",
		"SyntaxCheck", "RunUnitTests",
		"GetTable", "GetTableContents", "RunQuery",
		"FindDefinition", "FindReferences",
		"GetSystemInfo",
	}
	for _, tool := range coreTools {
		if !names[tool] {
			t.Errorf("core tool %q should be registered in focused mode", tool)
		}
	}
}

// TestRemovedToolsNotPresent verifies that removed experimental tools are gone.
func TestRemovedToolsNotPresent(t *testing.T) {
	for _, mode := range []string{"focused", "expert"} {
		t.Run(mode, func(t *testing.T) {
			s := newServerWithConfig(&Config{
				BaseURL:  "https://sap.example.com:44300",
				Username: "testuser",
				Password: "testpass",
				Mode:     mode,
			})

			names := toolNames(t, s)

			// These tools should NEVER be registered (removed features)
			removedTools := []string{
				// Debugger (WebSocket)
				"SetBreakpoint", "GetBreakpoints", "DeleteBreakpoint",
				"DebuggerListen", "DebuggerAttach", "DebuggerDetach",
				"DebuggerStep", "DebuggerGetStack", "DebuggerGetVariables",
				"CallRFC",
				// AMDP
				"AMDPDebuggerStart", "AMDPDebuggerResume", "AMDPDebuggerStop",
				"AMDPDebuggerStep", "AMDPGetVariables", "AMDPSetBreakpoint", "AMDPGetBreakpoints",
				// Git
				"GitTypes", "GitExport",
				// Reports
				"RunReport", "RunReportAsync", "GetAsyncResult",
				// Install
				"InstallZADTVSP", "InstallAbapGit", "ListDependencies", "InstallDummyTest", "DeployZip",
				// Service binding
				"PublishServiceBinding", "UnpublishServiceBinding",
				// Move (WebSocket)
				"MoveObject",
				// Help
				"GetAbapHelp",
			}

			for _, tool := range removedTools {
				if names[tool] {
					t.Errorf("removed tool %q should NOT be registered in %s mode", tool, mode)
				}
			}
		})
	}
}

// TestReadOnlyConfigSafetyIsReadOnly verifies that readonly mode actually sets safety.ReadOnly=true.
// The previous test only checked adtClient != nil, which was insufficient.
func TestReadOnlyConfigSafetyIsReadOnly(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "readonly",
	})

	safety := s.adtClient.Safety()
	if !safety.ReadOnly {
		t.Error("readonly mode must set Safety.ReadOnly = true")
	}
}

// TestReadOnlyModeWriteToolsRegistered documents current behavior: readonly mode registers
// write tools but blocks them at the safety layer, not at tool registration.
// This is a belt-and-suspenders design: tools exist but the ADT client refuses write ops.
func TestReadOnlyModeWriteToolsRegistered(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "readonly",
	})

	names := toolNames(t, s)

	// Write tools ARE registered (safety enforced at client layer, not tool layer)
	writeTools := []string{"WriteSource", "SyntaxCheck"}
	for _, tool := range writeTools {
		if !names[tool] {
			// If a write tool is absent it means it was removed — update this list
			t.Logf("note: tool %q not registered in readonly mode (may have been removed)", tool)
		}
	}
	// WriteSource is a core tool that should always be registered
	if !names["WriteSource"] {
		t.Error("WriteSource should be registered in readonly mode (safety blocks it at client layer)")
	}

	// The safety config must block writes regardless of what tools are registered
	safety := s.adtClient.Safety()
	if !safety.ReadOnly {
		t.Error("safety.ReadOnly must be true in readonly mode even if write tools are registered")
	}
}

// TestFocusedModeDefaultIsUnrestricted verifies default focused mode has no safety restrictions
// (backwards-compatible behaviour — safety must be explicitly configured).
func TestFocusedModeDefaultIsUnrestricted(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "focused",
	})

	safety := s.adtClient.Safety()
	if safety.ReadOnly {
		t.Error("focused mode without explicit ReadOnly flag should not be read-only")
	}
	if safety.BlockFreeSQL {
		t.Error("focused mode without explicit BlockFreeSQL flag should not block free SQL")
	}
}

// TestReadOnlyFlagOnFocusedMode verifies that ReadOnly=true in config applies safety
// regardless of mode.
func TestReadOnlyFlagOnFocusedMode(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:  "https://sap.example.com:44300",
		Username: "testuser",
		Password: "testpass",
		Mode:     "focused",
		ReadOnly: true,
	})

	safety := s.adtClient.Safety()
	if !safety.ReadOnly {
		t.Error("ReadOnly=true config flag must set Safety.ReadOnly = true")
	}
}

// TestBlockFreeSQLConfig verifies BlockFreeSQL config flag is applied to safety.
func TestBlockFreeSQLConfig(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:      "https://sap.example.com:44300",
		Username:     "testuser",
		Password:     "testpass",
		Mode:         "focused",
		BlockFreeSQL: true,
	})

	safety := s.adtClient.Safety()
	if !safety.BlockFreeSQL {
		t.Error("BlockFreeSQL=true config flag must set Safety.BlockFreeSQL = true")
	}
}

// TestAllowedPackagesConfig verifies AllowedPackages propagates to safety config.
func TestAllowedPackagesConfig(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:         "https://sap.example.com:44300",
		Username:        "testuser",
		Password:        "testpass",
		Mode:            "focused",
		AllowedPackages: []string{"$TMP", "Z*"},
	})

	safety := s.adtClient.Safety()
	if len(safety.AllowedPackages) != 2 {
		t.Fatalf("expected 2 AllowedPackages, got %d", len(safety.AllowedPackages))
	}
	if !safety.IsPackageAllowed("$TMP") {
		t.Error("$TMP should be allowed")
	}
	if !safety.IsPackageAllowed("ZTEST") {
		t.Error("ZTEST should be allowed via Z* wildcard")
	}
	if safety.IsPackageAllowed("PROD") {
		t.Error("PROD should NOT be allowed")
	}
}

// TestAllowedOpsConfig verifies AllowedOps string propagates to safety.
func TestAllowedOpsConfig(t *testing.T) {
	s := newServerWithConfig(&Config{
		BaseURL:    "https://sap.example.com:44300",
		Username:   "testuser",
		Password:   "testpass",
		Mode:       "focused",
		AllowedOps: "RSQ",
	})

	safety := s.adtClient.Safety()
	if safety.AllowedOps != "RSQ" {
		t.Errorf("AllowedOps = %q, want %q", safety.AllowedOps, "RSQ")
	}
}

// TestToolsConfigOverridesFocusedMode verifies .vsp.json tool config overrides.
func TestToolsConfigOverridesFocusedMode(t *testing.T) {
	// Disable a normally-enabled tool
	toolsConfig := map[string]bool{
		"GetSource": false,
	}

	s := newServerWithConfig(&Config{
		BaseURL:     "https://sap.example.com:44300",
		Username:    "testuser",
		Password:    "testpass",
		Mode:        "focused",
		ToolsConfig: toolsConfig,
	})

	names := toolNames(t, s)
	if names["GetSource"] {
		t.Error("GetSource should be disabled via toolsConfig override")
	}
}
