package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Checks whether the expected capability is specified in the capabilities.
func contains(expected string, values []string) bool {
	for _, v := range values {
		if v == expected {
			return true
		}
	}
	return false
}

func containsDevice(expected *Device, values []*Device) bool {
	for _, d := range values {
		if d.Path == expected.Path &&
			d.Permissions == expected.Permissions &&
			d.FileMode == expected.FileMode &&
			d.Major == expected.Major &&
			d.Minor == expected.Minor &&
			d.Type == expected.Type {
			return true
		}
	}
	return false
}

func loadConfig(name string) (*Config, error) {
	f, err := os.Open(filepath.Join("../sample_configs", name))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var container *Config
	if err := json.NewDecoder(f).Decode(&container); err != nil {
		return nil, err
	}

	// Check that a config doesn't contain extra fields
	var configMap, abstractMap map[string]interface{}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(f).Decode(&abstractMap); err != nil {
		return nil, err
	}

	configData, err := json.Marshal(&container)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(configData, &configMap); err != nil {
		return nil, err
	}

	for k := range configMap {
		delete(abstractMap, k)
	}

	if len(abstractMap) != 0 {
		return nil, fmt.Errorf("unknown fields: %s", abstractMap)
	}

	return container, nil
}

func TestRemoveNamespace(t *testing.T) {
	ns := Namespaces{
		{Type: NEWNET},
	}
	if !ns.Remove(NEWNET) {
		t.Fatal("NEWNET was not removed")
	}
	if len(ns) != 0 {
		t.Fatalf("namespaces should have 0 items but reports %d", len(ns))
	}
}

func TestHostUIDNoUSERNS(t *testing.T) {
	config := &Config{
		Namespaces: Namespaces{},
	}
	uid, err := config.HostUID()
	if err != nil {
		t.Fatal(err)
	}
	if uid != 0 {
		t.Fatalf("expected uid 0 with no USERNS but received %d", uid)
	}
}

func TestHostUIDWithUSERNS(t *testing.T) {
	config := &Config{
		Namespaces: Namespaces{{Type: NEWUSER}},
		UidMappings: []IDMap{
			{
				ContainerID: 0,
				HostID:      1000,
				Size:        1,
			},
		},
	}
	uid, err := config.HostUID()
	if err != nil {
		t.Fatal(err)
	}
	if uid != 1000 {
		t.Fatalf("expected uid 1000 with no USERNS but received %d", uid)
	}
}

func TestHostGIDNoUSERNS(t *testing.T) {
	config := &Config{
		Namespaces: Namespaces{},
	}
	uid, err := config.HostGID()
	if err != nil {
		t.Fatal(err)
	}
	if uid != 0 {
		t.Fatalf("expected gid 0 with no USERNS but received %d", uid)
	}
}

func TestHostGIDWithUSERNS(t *testing.T) {
	config := &Config{
		Namespaces: Namespaces{{Type: NEWUSER}},
		GidMappings: []IDMap{
			{
				ContainerID: 0,
				HostID:      1000,
				Size:        1,
			},
		},
	}
	uid, err := config.HostGID()
	if err != nil {
		t.Fatal(err)
	}
	if uid != 1000 {
		t.Fatalf("expected gid 1000 with no USERNS but received %d", uid)
	}
}
