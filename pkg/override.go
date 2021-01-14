package rein

import (
	"encoding/json"
	"errors"
)

// Override is a Toggle value to be consumed by and restricted to a semantic version range
type Override struct {
	MinimumVersionStruct *Version    `json:"-"`
	MaximumVersionStruct *Version    `json:"-"`
	MinimumVersion       string      `json:"minimumVersion"`
	MaximumVersion       string      `json:"maximumVersion"`
	Value                interface{} `json:"value"`
}

// UnmarshalJSON Custom UnmarshalJSON method for validating Override
func (o *Override) UnmarshalJSON(b []byte) error {

	var alias overrideAlias

	err := json.Unmarshal(b, &alias)
	if err != nil {
		return err
	}

	*o = alias.toOverride()

	minEmpty := o.MinimumVersion == ""
	maxEmpty := o.MaximumVersion == ""

	if minEmpty && maxEmpty {
		return errors.New("Override ranges cannot both be empty")
	}

	o.MinimumVersionStruct, err = NewVersion(o.MinimumVersion)
	if err != nil {
		return err
	}
	o.MaximumVersionStruct, err = NewVersion(o.MaximumVersion)
	if err != nil {
		return err
	}

	return nil
}

type overrideAlias Override

func (o overrideAlias) toOverride() Override {
	return Override{
		nil,
		nil,
		o.MinimumVersion,
		o.MaximumVersion,
		o.Value,
	}
}
