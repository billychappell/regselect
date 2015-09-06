// +build windows
package regselect

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/sys/windows/registry"
)

// Property contains the name, type, value and previous value located within the
// current key. The registry is updated using the Value field after recording its
// pre-op value in the PrevValue field whenever the Set() method is called.
type Property struct {
	Name      string
	Type      string
	Value     interface{}
	PrevValue interface{}
}

// Key uses the Path and Scope fields when opening a particular key from the registry.
// Properties lists all the values you want to change within a given key.
type Key struct {
	Path, Scope string
	Properties  []Property
}

// GetScope() returns a scope's corresponding const value for use with
// the 'golang.org/x/sys/windows/registry' library.
func (k *Key) GetScope() registry.Key {
	switch k.Scope {
	case "LOCAL_MACHINE":
		return registry.LOCAL_MACHINE
	case "CURRENT_USER":
		return registry.CURRENT_USER
	case "CLASSES_ROOT":
		return registry.CLASSES_ROOT
	case "CURRENT_CONFIG":
		return registry.CURRENT_CONFIG
	case "USERS":
		return registry.USERS
	default:
		return registry.LOCAL_MACHINE
	}
}

// Config is the base struct that we use to Unmarshal JSON config files
// into a Go-readable data structure.
type Config []Key

// Validate() is a method implemented by the Config struct that attempts to
// retrieve the value of each property for each key in a Config. Each Property's
// PrevValue field is updated to its current value. This should be called before
// Set() but not after Write() is implemented.
func (c *Config) Validate() (err error) {
	con := *c
	for i, v := range con {
		k, err := registry.OpenKey(v.GetScope(), v.Path, registry.ALL_ACCESS)
		if err != nil {
			return err
		}
		defer k.Close()
		for n, p := range v.Properties {
			val := &con[i].Properties[n].PrevValue
			switch p.Type {
			case "DWord", "QWord":
				if s, _, err := k.GetIntegerValue(p.Name); err == nil {
					*val = s
				} else {
					return err
				}
			case "String":
				if s, _, err := k.GetStringValue(p.Name); err == nil {
					*val = s
				} else {
					return err
				}
			case "Strings":
				if s, _, err := k.GetStringsValue(p.Name); err == nil {
					*val = s
				} else {
					return err
				}
			case "Binary":
				if s, _, err := k.GetBinaryValue(p.Name); err == nil {
					*val = s
				} else {
					return err
				}
			default:
				var buf []byte
				if _, _, err := k.GetValue(p.Name, buf); err != nil {
					return err
				}
				return fmt.Errorf("%s of %s path in %s scope returned code %d.") // TODO: Convert const int representation of value types to explicitly match what the user should type into their JSON config.
			}
		}
	}
	return nil
}

// Set() writes the JSON config to the registry. Make sure you call Validate() before calling Set() to
// avoid filling values with the wrong type. For now, errors are handled by stopping mid-operation. Looking
// at possibly reverting all values back to prev values or providing more detailed data to users/letting them
// choose.
func (c *Config) Set() (err error) {
	for _, v := range *c {
		k, err := registry.OpenKey(v.GetScope(), v.Path, registry.ALL_ACCESS)
		if err != nil {
			return err
		}
		defer k.Close()

		for _, p := range v.Properties {
			switch p.Type {
			case "DWord":
				if err := k.SetDWordValue(p.Name, uint32(p.Value.(float64))); err != nil {
					return err
				}
			case "QWord":
				if err := k.SetQWordValue(p.Name, uint64(p.Value.(float64))); err != nil {
					return err
				}
			case "String":
				if err := k.SetStringValue(p.Name, p.Value.(string)); err != nil {
					return err
				}
			case "Strings":
				if err := k.SetStringsValue(p.Name, p.Value.([]string)); err != nil {
					return err
				}
			case "Binary":
				if err := k.SetBinaryValue(p.Name, p.Value.([]byte)); err != nil {
					return err
				}
			default:
				return fmt.Errorf("please check the type of %s in %s is correctly set. Currently: %s", p.Name, v.Path, p.Type)
			}
		}
	}
	return nil
}

// Write() outputs the current Config to filename in JSON format.
// Usually, you should call this method anytime you use Set() to
// save your pre-op registry values in the PrevValue field equivalent.
func (c *Config) Write(filename string) (err error) {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	if _, err := file.Write(b); err != nil {
		return err
	}

	return nil
}

// Unmarshal accepts a JSON filename as an argument and returns a Config struct
// to use for editing values from the registry.
func Unmarshal(filename string) (c *Config, err error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}
	if err = json.Unmarshal(dat, &c); err != nil {
		return c, err
	}

	return c, nil
}
