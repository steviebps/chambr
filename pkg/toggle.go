package realm

import (
	"encoding/json"
	"fmt"
	"reflect"

	"golang.org/x/mod/semver"
)

type ToggleType string

const (
	booleanType ToggleType = "boolean"
	stringType  ToggleType = "string"
	numberType  ToggleType = "number"
	customType  ToggleType = "custom"
)

// Toggle is a feature switch/toggle structure for holding
// its name, value, type and any overrides to be parsed by the applicable realm sdk
type Toggle struct {
	Name      string      `json:"name"`
	Type      ToggleType  `json:"type"`
	Value     interface{} `json:"value"`
	Overrides []*Override `json:"overrides,omitempty"`
	// ToggleValidator ToggleValidator `json:"-"`
}

// type ToggleValidator interface {
// 	ValidateValue(value interface{}) bool
// 	GetValueAt(version string) interface{}
// }

type toggleAlias Toggle

func (t toggleAlias) toToggle() Toggle {
	return Toggle(t)
}

// IsValidValue determines whether or not the passed value's type matches the ToggleType
func (t Toggle) IsValidValue(value interface{}) bool {
	typeOfValue := reflect.TypeOf(value).String()

	switch typeOfValue {
	case "bool":
		return t.Type == booleanType
	case "string":
		return t.Type == stringType
	case "float64":
		return t.Type == numberType
	default:
		return false
	}
}

// UnmarshalJSON Custom UnmarshalJSON method for validating toggle Value to the ToggleType
func (t *Toggle) UnmarshalJSON(b []byte) error {
	var alias toggleAlias

	err := json.Unmarshal(b, &alias)
	if err != nil {
		return err
	}
	*t = alias.toToggle()

	if !t.IsValidValue(t.Value) {
		return fmt.Errorf("%v (%T) not of the type %q from the toggle: %s", t.Value, t.Value, t.Type, t.Name)
	}

	var previous *Override
	for _, override := range t.Overrides {
		// overrides should not overlap
		if previous != nil && semver.Compare(previous.MaximumVersion, override.MinimumVersion) == 1 {
			return fmt.Errorf("an override with maximum version %v is semantically greater than the next override's minimum version (%v) ", previous.MaximumVersion, override.MinimumVersion)
		}

		if !t.IsValidValue(override.Value) {
			return fmt.Errorf("%v (%T) not of the type %q from the toggle override: %s", override.Value, override.Value, t.Type, t.Name)
		}

		previous = override
	}

	return nil
}

// GetValueAt returns the value at the given version.
// Will return default value if version is empty string or no override is present for the specified version
func (t *Toggle) GetValueAt(version string) interface{} {
	if version != "" {
		if override := t.GetOverride(version); override != nil {
			return override.Value
		}
	}

	return t.Value
}

// GetOverride returns the first override that encapsulates version passed
func (t *Toggle) GetOverride(version string) *Override {

	for _, override := range t.Overrides {
		if semver.Compare(override.MinimumVersion, version) <= 0 && semver.Compare(override.MaximumVersion, version) >= 0 {
			return override
		}
	}

	return nil
}
