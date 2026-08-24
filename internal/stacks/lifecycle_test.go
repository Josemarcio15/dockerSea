package stacks

import (
	"strings"
	"testing"
)

func TestGenerateOperationID(t *testing.T) {
	op1 := GenerateOperationID()
	op2 := GenerateOperationID()

	if !strings.HasPrefix(op1, "op_") {
		t.Fatalf("esperado prefixo 'op_', obtido %s", op1)
	}
	if op1 == op2 {
		t.Fatalf("IDs de operação consecutivos devem ser únicos: %s == %s", op1, op2)
	}
}

func TestBuildRemoteLifecycleDir(t *testing.T) {
	tests := []struct {
		stackID     string
		operationID string
		expectErr   bool
		expected    string
	}{
		{
			stackID:     "stk_123",
			operationID: "op_456",
			expectErr:   false,
			expected:    "$HOME/.docksea/lifecycle/stk_123/op_456",
		},
		{
			stackID:     "stk-valid",
			operationID: "op-valid",
			expectErr:   false,
			expected:    "$HOME/.docksea/lifecycle/stk-valid/op-valid",
		},
		{
			stackID:     "../malicious",
			operationID: "op_1",
			expectErr:   true,
		},
		{
			stackID:     "stk_1",
			operationID: "../../escape",
			expectErr:   true,
		},
		{
			stackID:     "",
			operationID: "op_1",
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		res, err := BuildRemoteLifecycleDir(tt.stackID, tt.operationID)
		if tt.expectErr {
			if err == nil {
				t.Errorf("esperado erro para stackID=%s, operationID=%s", tt.stackID, tt.operationID)
			}
		} else {
			if err != nil {
				t.Errorf("erro inesperado: %v", err)
			}
			if res != tt.expected {
				t.Errorf("esperado %s, obtido %s", tt.expected, res)
			}
		}
	}
}
