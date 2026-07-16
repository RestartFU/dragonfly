package world

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func TestNewCustomBlockRegistryRejectsExcessStates(t *testing.T) {
	properties := make([]any, 17)
	for i := range properties {
		properties[i] = map[string]any{
			"name": fmt.Sprintf("server:property_%d", i),
			"enum": []any{false, true},
		}
	}
	_, err := NewCustomBlockRegistry([]protocol.BlockEntry{{
		Name:       "server:test",
		Properties: map[string]any{"properties": properties},
	}})
	if err == nil || !strings.Contains(err.Error(), "exceed limit") {
		t.Fatalf("build registry error = %v, want state limit error", err)
	}
}

func TestNewCustomBlockRegistryMarksAppendedStatesNBT(t *testing.T) {
	registry, err := NewCustomBlockRegistry([]protocol.BlockEntry{{Name: "server:test"}})
	if err != nil {
		t.Fatalf("build registry: %v", err)
	}
	runtimeID, ok := registry.StateToRuntimeID("server:test", nil)
	if !ok {
		t.Fatal("custom state missing from registry")
	}
	if !registry.NBTBlock(runtimeID) {
		t.Fatal("custom state is not marked NBT-capable")
	}
}
