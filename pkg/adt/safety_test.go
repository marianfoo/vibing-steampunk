package adt

import (
	"testing"
)

func TestSafetyConfig_IsOperationAllowed(t *testing.T) {
	tests := []struct {
		name     string
		config   SafetyConfig
		op       OperationType
		expected bool
	}{
		{
			name:     "ReadOnly blocks create",
			config:   SafetyConfig{ReadOnly: true},
			op:       OpCreate,
			expected: false,
		},
		{
			name:     "ReadOnly blocks update",
			config:   SafetyConfig{ReadOnly: true},
			op:       OpUpdate,
			expected: false,
		},
		{
			name:     "ReadOnly blocks delete",
			config:   SafetyConfig{ReadOnly: true},
			op:       OpDelete,
			expected: false,
		},
		{
			name:     "ReadOnly blocks activate",
			config:   SafetyConfig{ReadOnly: true},
			op:       OpActivate,
			expected: false,
		},
		{
			name:     "ReadOnly blocks workflow",
			config:   SafetyConfig{ReadOnly: true},
			op:       OpWorkflow,
			expected: false,
		},
		{
			name:     "ReadOnly allows read",
			config:   SafetyConfig{ReadOnly: true},
			op:       OpRead,
			expected: true,
		},
		{
			name:     "ReadOnly allows search",
			config:   SafetyConfig{ReadOnly: true},
			op:       OpSearch,
			expected: true,
		},
		{
			name:     "ReadOnly allows query",
			config:   SafetyConfig{ReadOnly: true},
			op:       OpQuery,
			expected: true,
		},
		{
			name:     "BlockFreeSQL blocks free SQL",
			config:   SafetyConfig{BlockFreeSQL: true},
			op:       OpFreeSQL,
			expected: false,
		},
		{
			name:     "BlockFreeSQL allows read",
			config:   SafetyConfig{BlockFreeSQL: true},
			op:       OpRead,
			expected: true,
		},
		{
			name:     "AllowedOps whitelist - allowed op",
			config:   SafetyConfig{AllowedOps: "RSQ"},
			op:       OpRead,
			expected: true,
		},
		{
			name:     "AllowedOps whitelist - blocked op",
			config:   SafetyConfig{AllowedOps: "RSQ"},
			op:       OpCreate,
			expected: false,
		},
		{
			name:     "DisallowedOps blacklist",
			config:   SafetyConfig{DisallowedOps: "CDUA"},
			op:       OpCreate,
			expected: false,
		},
		{
			name:     "DisallowedOps allows non-blacklisted",
			config:   SafetyConfig{DisallowedOps: "CDUA"},
			op:       OpRead,
			expected: true,
		},
		{
			name:     "DisallowedOps overrides AllowedOps",
			config:   SafetyConfig{AllowedOps: "RCDU", DisallowedOps: "CD"},
			op:       OpCreate,
			expected: false,
		},
		{
			name:     "DryRun allows all",
			config:   SafetyConfig{DryRun: true, ReadOnly: true},
			op:       OpCreate,
			expected: true,
		},
		{
			name:     "Unrestricted allows all",
			config:   UnrestrictedSafetyConfig(),
			op:       OpCreate,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsOperationAllowed(tt.op)
			if result != tt.expected {
				t.Errorf("IsOperationAllowed(%c) = %v, expected %v", tt.op, result, tt.expected)
			}
		})
	}
}

func TestSafetyConfig_CheckOperation(t *testing.T) {
	config := SafetyConfig{ReadOnly: true}

	// Should allow read
	err := config.CheckOperation(OpRead, "GetClass")
	if err != nil {
		t.Errorf("CheckOperation(OpRead) should not error, got: %v", err)
	}

	// Should block create
	err = config.CheckOperation(OpCreate, "CreateObject")
	if err == nil {
		t.Error("CheckOperation(OpCreate) should error in read-only mode")
	}
}

func TestSafetyConfig_IsPackageAllowed(t *testing.T) {
	tests := []struct {
		name     string
		config   SafetyConfig
		pkg      string
		expected bool
	}{
		{
			name:     "Empty AllowedPackages allows all",
			config:   SafetyConfig{},
			pkg:      "ZANY",
			expected: true,
		},
		{
			name:     "Exact match",
			config:   SafetyConfig{AllowedPackages: []string{"$TMP", "ZTEST"}},
			pkg:      "$TMP",
			expected: true,
		},
		{
			name:     "Not in whitelist",
			config:   SafetyConfig{AllowedPackages: []string{"$TMP", "ZTEST"}},
			pkg:      "ZPROD",
			expected: false,
		},
		{
			name:     "Wildcard match - Z*",
			config:   SafetyConfig{AllowedPackages: []string{"Z*"}},
			pkg:      "ZTEST",
			expected: true,
		},
		{
			name:     "Wildcard match - $*",
			config:   SafetyConfig{AllowedPackages: []string{"$*"}},
			pkg:      "$TMP",
			expected: true,
		},
		{
			name:     "Wildcard no match",
			config:   SafetyConfig{AllowedPackages: []string{"Z*"}},
			pkg:      "$TMP",
			expected: false,
		},
		{
			name:     "Case insensitive",
			config:   SafetyConfig{AllowedPackages: []string{"ztest"}},
			pkg:      "ZTEST",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsPackageAllowed(tt.pkg)
			if result != tt.expected {
				t.Errorf("IsPackageAllowed(%s) = %v, expected %v", tt.pkg, result, tt.expected)
			}
		})
	}
}

func TestSafetyConfig_CheckPackage(t *testing.T) {
	config := SafetyConfig{AllowedPackages: []string{"$TMP", "Z*"}}

	// Should allow $TMP
	err := config.CheckPackage("$TMP")
	if err != nil {
		t.Errorf("CheckPackage($TMP) should not error, got: %v", err)
	}

	// Should allow ZTEST (wildcard)
	err = config.CheckPackage("ZTEST")
	if err != nil {
		t.Errorf("CheckPackage(ZTEST) should not error, got: %v", err)
	}

	// Should block PROD
	err = config.CheckPackage("PROD")
	if err == nil {
		t.Error("CheckPackage(PROD) should error")
	}
}

func TestSafetyConfig_String(t *testing.T) {
	tests := []struct {
		name     string
		config   SafetyConfig
		contains []string
	}{
		{
			name:     "Unrestricted",
			config:   UnrestrictedSafetyConfig(),
			contains: []string{"UNRESTRICTED"},
		},
		{
			name:     "ReadOnly",
			config:   SafetyConfig{ReadOnly: true},
			contains: []string{"READ-ONLY"},
		},
		{
			name:     "BlockFreeSQL",
			config:   SafetyConfig{BlockFreeSQL: true},
			contains: []string{"NO-FREE-SQL"},
		},
		{
			name:     "DryRun",
			config:   SafetyConfig{DryRun: true},
			contains: []string{"DRY-RUN"},
		},
		{
			name:     "AllowedOps",
			config:   SafetyConfig{AllowedOps: "RSQ"},
			contains: []string{"AllowedOps=RSQ"},
		},
		{
			name:     "Multiple flags",
			config:   SafetyConfig{ReadOnly: true, BlockFreeSQL: true},
			contains: []string{"READ-ONLY", "NO-FREE-SQL"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.String()
			for _, expected := range tt.contains {
				if !contains(result, expected) {
					t.Errorf("String() = %q, should contain %q", result, expected)
				}
			}
		})
	}
}

func TestDefaultSafetyConfig(t *testing.T) {
	config := DefaultSafetyConfig()

	if !config.ReadOnly {
		t.Error("DefaultSafetyConfig should be read-only")
	}

	if !config.BlockFreeSQL {
		t.Error("DefaultSafetyConfig should block free SQL")
	}

	if !config.IsOperationAllowed(OpRead) {
		t.Error("DefaultSafetyConfig should allow read operations")
	}

	if config.IsOperationAllowed(OpCreate) {
		t.Error("DefaultSafetyConfig should not allow create operations")
	}
}

func TestDevelopmentSafetyConfig(t *testing.T) {
	config := DevelopmentSafetyConfig()

	if config.ReadOnly {
		t.Error("DevelopmentSafetyConfig should not be read-only")
	}

	if !config.BlockFreeSQL {
		t.Error("DevelopmentSafetyConfig should block free SQL")
	}

	if !config.IsPackageAllowed("$TMP") {
		t.Error("DevelopmentSafetyConfig should allow $TMP")
	}

	if config.IsPackageAllowed("ZPROD") {
		t.Error("DevelopmentSafetyConfig should not allow ZPROD")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInMiddle(s, substr)))
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestSafetyConfig_CheckTransportableEdit(t *testing.T) {
	tests := []struct {
		name        string
		config      SafetyConfig
		transport   string
		opName      string
		expectError bool
	}{
		{
			name:        "No transport - always allowed",
			config:      SafetyConfig{},
			transport:   "",
			opName:      "EditSource",
			expectError: false,
		},
		{
			name:        "Transport provided but not allowed - blocked",
			config:      SafetyConfig{AllowTransportableEdits: false},
			transport:   "DEVK900123",
			opName:      "EditSource",
			expectError: true,
		},
		{
			name:        "Transport provided and allowed - success",
			config:      SafetyConfig{AllowTransportableEdits: true},
			transport:   "DEVK900123",
			opName:      "EditSource",
			expectError: false,
		},
		{
			name:        "Transport allowed but not in whitelist - blocked",
			config:      SafetyConfig{AllowTransportableEdits: true, AllowedTransports: []string{"A4HK*"}},
			transport:   "DEVK900123",
			opName:      "WriteSource",
			expectError: true,
		},
		{
			name:        "Transport allowed and in whitelist - success",
			config:      SafetyConfig{AllowTransportableEdits: true, AllowedTransports: []string{"DEVK*"}},
			transport:   "DEVK900123",
			opName:      "WriteSource",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.CheckTransportableEdit(tt.transport, tt.opName)
			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestSafetyConfig_CheckTransportableEdit_ErrorMessage(t *testing.T) {
	config := SafetyConfig{AllowTransportableEdits: false}
	err := config.CheckTransportableEdit("DEVK900123", "EditSource")

	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	errMsg := err.Error()

	// Check that error message contains helpful information
	if !contains(errMsg, "EditSource") {
		t.Error("Error message should contain operation name")
	}
	if !contains(errMsg, "DEVK900123") {
		t.Error("Error message should contain transport number")
	}
	if !contains(errMsg, "--allow-transportable-edits") {
		t.Error("Error message should mention CLI flag")
	}
	if !contains(errMsg, "SAP_ALLOW_TRANSPORTABLE_EDITS") {
		t.Error("Error message should mention environment variable")
	}
}

// TestSafetyConfig_IsPackageAllowed_EmptyString documents the edge case from issue #71.
// When AllowedPackages is set, an empty package name (common for top-level packages) is NOT
// in the whitelist. The CreateObject function must NOT call CheckPackage with PackageName
// when creating packages — it must use the package Name itself.
func TestSafetyConfig_IsPackageAllowed_EmptyString(t *testing.T) {
	config := SafetyConfig{AllowedPackages: []string{"$TMP", "Z*"}}

	// Empty string is not in the whitelist — returns false
	if config.IsPackageAllowed("") {
		t.Error("empty package name should NOT be allowed when whitelist is non-empty (regression: issue #71)")
	}

	// Verify the real fix: the package being created ($ZTEST) IS allowed
	if !config.IsPackageAllowed("$TMP") {
		t.Error("$TMP should be allowed")
	}
}

// TestSafetyConfig_CheckOperation_AllTypes verifies all operation types behave correctly
// under ReadOnly mode.
func TestSafetyConfig_CheckOperation_AllTypes(t *testing.T) {
	readonly := SafetyConfig{ReadOnly: true}

	blockedOps := []struct {
		op   OperationType
		name string
	}{
		{OpCreate, "OpCreate"},
		{OpUpdate, "OpUpdate"},
		{OpDelete, "OpDelete"},
		{OpActivate, "OpActivate"},
		{OpWorkflow, "OpWorkflow"},
	}

	for _, tc := range blockedOps {
		if readonly.IsOperationAllowed(tc.op) {
			t.Errorf("ReadOnly should block %s", tc.name)
		}
	}

	allowedOps := []struct {
		op   OperationType
		name string
	}{
		{OpRead, "OpRead"},
		{OpSearch, "OpSearch"},
		{OpQuery, "OpQuery"},
		{OpTest, "OpTest"},
		{OpLock, "OpLock"},
		{OpIntelligence, "OpIntelligence"},
	}

	for _, tc := range allowedOps {
		if !readonly.IsOperationAllowed(tc.op) {
			t.Errorf("ReadOnly should allow %s", tc.name)
		}
	}
}

// TestSafetyConfig_DefaultConfig_AllowedOps verifies the default config allows RSQTI ops
// and blocks write ops.
func TestSafetyConfig_DefaultConfig_AllowedOps(t *testing.T) {
	config := DefaultSafetyConfig()

	allowedByDefault := []struct {
		op   OperationType
		name string
	}{
		{OpRead, "OpRead"},
		{OpSearch, "OpSearch"},
		{OpQuery, "OpQuery"},
		{OpTest, "OpTest"},
		{OpIntelligence, "OpIntelligence"},
	}

	for _, tc := range allowedByDefault {
		if !config.IsOperationAllowed(tc.op) {
			t.Errorf("DefaultSafetyConfig should allow %s", tc.name)
		}
	}

	blockedByDefault := []struct {
		op   OperationType
		name string
	}{
		{OpCreate, "OpCreate"},
		{OpUpdate, "OpUpdate"},
		{OpDelete, "OpDelete"},
		{OpActivate, "OpActivate"},
		{OpWorkflow, "OpWorkflow"},
		{OpFreeSQL, "OpFreeSQL"},
		{OpTransport, "OpTransport"},
	}

	for _, tc := range blockedByDefault {
		if config.IsOperationAllowed(tc.op) {
			t.Errorf("DefaultSafetyConfig should block %s", tc.name)
		}
	}
}

// TestSafetyConfig_IsTransportAllowed tests the IsTransportAllowed method.
func TestSafetyConfig_IsTransportAllowed(t *testing.T) {
	tests := []struct {
		name      string
		config    SafetyConfig
		transport string
		expected  bool
	}{
		{
			name:      "Transports disabled - always false",
			config:    SafetyConfig{EnableTransports: false},
			transport: "DEVK900001",
			expected:  false,
		},
		{
			name:      "Transports enabled, no whitelist - always true",
			config:    SafetyConfig{EnableTransports: true},
			transport: "DEVK900001",
			expected:  true,
		},
		{
			name:      "Transports enabled, exact whitelist match",
			config:    SafetyConfig{EnableTransports: true, AllowedTransports: []string{"DEVK900001"}},
			transport: "DEVK900001",
			expected:  true,
		},
		{
			name:      "Transports enabled, not in whitelist",
			config:    SafetyConfig{EnableTransports: true, AllowedTransports: []string{"A4HK*"}},
			transport: "DEVK900001",
			expected:  false,
		},
		{
			name:      "Transports enabled, wildcard match",
			config:    SafetyConfig{EnableTransports: true, AllowedTransports: []string{"DEVK*"}},
			transport: "DEVK900001",
			expected:  true,
		},
		{
			name:      "Transports enabled, case insensitive",
			config:    SafetyConfig{EnableTransports: true, AllowedTransports: []string{"devk*"}},
			transport: "DEVK900001",
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsTransportAllowed(tt.transport)
			if result != tt.expected {
				t.Errorf("IsTransportAllowed(%q) = %v, expected %v", tt.transport, result, tt.expected)
			}
		})
	}
}

// TestSafetyConfig_IsTransportWriteAllowed verifies write permission checks on transports.
func TestSafetyConfig_IsTransportWriteAllowed(t *testing.T) {
	tests := []struct {
		name     string
		config   SafetyConfig
		expected bool
	}{
		{
			name:     "Transports disabled - no writes",
			config:   SafetyConfig{EnableTransports: false},
			expected: false,
		},
		{
			name:     "Transports enabled, not read-only - writes allowed",
			config:   SafetyConfig{EnableTransports: true, TransportReadOnly: false},
			expected: true,
		},
		{
			name:     "Transports enabled, read-only - no writes",
			config:   SafetyConfig{EnableTransports: true, TransportReadOnly: true},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsTransportWriteAllowed()
			if result != tt.expected {
				t.Errorf("IsTransportWriteAllowed() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestSafetyConfig_CheckTransport covers all combinations of read/write and transport flags.
func TestSafetyConfig_CheckTransport(t *testing.T) {
	tests := []struct {
		name        string
		config      SafetyConfig
		transport   string
		opName      string
		isWrite     bool
		expectError bool
	}{
		// Read operations
		{
			name:        "Read op, transports disabled, no AllowTransportableEdits - blocked",
			config:      SafetyConfig{EnableTransports: false, AllowTransportableEdits: false},
			transport:   "DEVK900001",
			opName:      "ListTransports",
			isWrite:     false,
			expectError: true,
		},
		{
			name:        "Read op, transports enabled - allowed",
			config:      SafetyConfig{EnableTransports: true},
			transport:   "DEVK900001",
			opName:      "GetTransport",
			isWrite:     false,
			expectError: false,
		},
		{
			name:        "Read op, AllowTransportableEdits=true (no EnableTransports) - allowed",
			config:      SafetyConfig{EnableTransports: false, AllowTransportableEdits: true},
			transport:   "DEVK900001",
			opName:      "ListTransports",
			isWrite:     false,
			expectError: false,
		},
		{
			name:        "Read op, transports enabled, not in whitelist - blocked",
			config:      SafetyConfig{EnableTransports: true, AllowedTransports: []string{"A4HK*"}},
			transport:   "DEVK900001",
			opName:      "GetTransport",
			isWrite:     false,
			expectError: true,
		},
		// Write operations
		{
			name:        "Write op, transports disabled - blocked",
			config:      SafetyConfig{EnableTransports: false},
			transport:   "DEVK900001",
			opName:      "CreateTransport",
			isWrite:     true,
			expectError: true,
		},
		{
			name:        "Write op, AllowTransportableEdits=true but no EnableTransports - blocked",
			config:      SafetyConfig{EnableTransports: false, AllowTransportableEdits: true},
			transport:   "DEVK900001",
			opName:      "CreateTransport",
			isWrite:     true,
			expectError: true,
		},
		{
			name:        "Write op, transports enabled - allowed",
			config:      SafetyConfig{EnableTransports: true},
			transport:   "DEVK900001",
			opName:      "CreateTransport",
			isWrite:     true,
			expectError: false,
		},
		{
			name:        "Write op, transports enabled but read-only - blocked",
			config:      SafetyConfig{EnableTransports: true, TransportReadOnly: true},
			transport:   "DEVK900001",
			opName:      "CreateTransport",
			isWrite:     true,
			expectError: true,
		},
		{
			name:        "Write op, transports enabled, in whitelist - allowed",
			config:      SafetyConfig{EnableTransports: true, AllowedTransports: []string{"DEVK*"}},
			transport:   "DEVK900001",
			opName:      "CreateTransport",
			isWrite:     true,
			expectError: false,
		},
		{
			name:        "Write op, transports enabled, not in whitelist - blocked",
			config:      SafetyConfig{EnableTransports: true, AllowedTransports: []string{"A4HK*"}},
			transport:   "DEVK900001",
			opName:      "CreateTransport",
			isWrite:     true,
			expectError: true,
		},
		// List operation (empty transport = wildcard)
		{
			name:        "List op (empty transport), transports enabled - allowed",
			config:      SafetyConfig{EnableTransports: true},
			transport:   "",
			opName:      "ListTransports",
			isWrite:     false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.CheckTransport(tt.transport, tt.opName, tt.isWrite)
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// TestSafetyConfig_CheckTransport_ErrorMessages verifies that error messages are helpful.
func TestSafetyConfig_CheckTransport_ErrorMessages(t *testing.T) {
	// Write op with AllowTransportableEdits but no EnableTransports — specific error
	config := SafetyConfig{EnableTransports: false, AllowTransportableEdits: true}
	err := config.CheckTransport("DEVK900001", "CreateTransport", true)
	if err == nil {
		t.Fatal("Expected error")
	}
	if !contains(err.Error(), "--enable-transports") {
		t.Errorf("Error should suggest --enable-transports flag, got: %v", err)
	}

	// No EnableTransports at all — generic error
	config2 := SafetyConfig{EnableTransports: false}
	err2 := config2.CheckTransport("DEVK900001", "ListTransports", false)
	if err2 == nil {
		t.Fatal("Expected error")
	}
	if !contains(err2.Error(), "SAP_ENABLE_TRANSPORTS") {
		t.Errorf("Error should mention SAP_ENABLE_TRANSPORTS, got: %v", err2)
	}
}

// TestSafetyConfig_OpTransport_RequiresEnableTransports verifies transport operations
// always require explicit opt-in via IsOperationAllowed.
func TestSafetyConfig_OpTransport_RequiresEnableTransports(t *testing.T) {
	// Without EnableTransports, OpTransport is always blocked
	noTransport := SafetyConfig{}
	if noTransport.IsOperationAllowed(OpTransport) {
		t.Error("OpTransport should be blocked when EnableTransports is false")
	}

	// Even with unrestricted config, transport requires explicit EnableTransports
	unrestricted := UnrestrictedSafetyConfig()
	if unrestricted.IsOperationAllowed(OpTransport) {
		t.Error("OpTransport should be blocked even in unrestricted config without EnableTransports")
	}

	// With EnableTransports set, OpTransport is allowed
	withTransport := SafetyConfig{EnableTransports: true}
	if !withTransport.IsOperationAllowed(OpTransport) {
		t.Error("OpTransport should be allowed when EnableTransports is true")
	}
}

// TestSafetyConfig_String_TransportSettings verifies transport-related flags in String output.
func TestSafetyConfig_String_TransportSettings(t *testing.T) {
	config := SafetyConfig{
		EnableTransports:        true,
		TransportReadOnly:       true,
		AllowedTransports:       []string{"DEVK*"},
		AllowTransportableEdits: true,
	}

	result := config.String()

	expected := []string{"TRANSPORTS-ENABLED", "TRANSPORT-READ-ONLY", "AllowedTransports=", "TRANSPORTABLE-EDITS-ALLOWED"}
	for _, s := range expected {
		if !contains(result, s) {
			t.Errorf("String() should contain %q, got: %q", s, result)
		}
	}
}

// TestSafetyConfig_AllowedPackages_NoRestriction verifies zero-length whitelist allows all.
func TestSafetyConfig_AllowedPackages_NoRestriction(t *testing.T) {
	config := SafetyConfig{}

	packages := []string{"$TMP", "ZPROD", "SAP_BASIS", "", "ANYTHING"}
	for _, pkg := range packages {
		if !config.IsPackageAllowed(pkg) {
			t.Errorf("IsPackageAllowed(%q) = false, should be true when AllowedPackages is empty", pkg)
		}
	}
}

// TestSafetyConfig_DisallowedOps_EmptyAllowsAll verifies empty DisallowedOps blocks nothing.
func TestSafetyConfig_DisallowedOps_EmptyAllowsAll(t *testing.T) {
	config := SafetyConfig{DisallowedOps: ""}

	ops := []OperationType{OpRead, OpSearch, OpCreate, OpUpdate, OpDelete, OpActivate}
	for _, op := range ops {
		if !config.IsOperationAllowed(op) {
			t.Errorf("Empty DisallowedOps should not block op %c", op)
		}
	}
}

// TestSafetyConfig_AllowedOps_EmptyAllowsAll verifies empty AllowedOps allows everything.
func TestSafetyConfig_AllowedOps_EmptyAllowsAll(t *testing.T) {
	config := SafetyConfig{AllowedOps: ""}

	ops := []OperationType{OpRead, OpSearch, OpCreate, OpUpdate, OpDelete, OpActivate}
	for _, op := range ops {
		if !config.IsOperationAllowed(op) {
			t.Errorf("Empty AllowedOps should allow op %c", op)
		}
	}
}
