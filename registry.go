package regselect

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/sys/windows/registry"
)

// Property holds the name, type and value of a registry variable that you wish
// to configure. Once you have changed the value of a Property with this package,
// the previous value is also stored  .
type Property struct {
	Name      string
	Type      ValType
	Value     interface{}
	PrevValue interface{}
}

// ValType is used to determine which method to use when getting or changing the
// value of a Property.
type ValType string

// Key stores related properties that share the same scope and path.
type Key struct {
	Path       string
	Scope      string
	Properties []Property
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

// Unmarshal accepts a JSON filename as an argument and returns a Config struct
// to use for editing values from the registry.
func Unmarshal(filename string) (c Config, err error) {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(dat, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// Write creates a JSON file using a Config. Filename is passed as the only
// argument. This method is particularly useful for updating a config file with
// each Property's prevValue after a recent change.
func (c Config) Write(filename string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Unable to marshal into JSON: \n %v", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Unable to create file, %s: \n %v", filename, err)
	}

	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("Unable to write to file %s: \n %v", filename, err)
	}

	return nil
}

func (p *Property) getError(err error) error {
	return fmt.Errorf("can't get %s value for %s: \n %v", p.Type, p.Name, err)
}

func (p *Property) setError(err error) error {
	return fmt.Errorf("can't set %s value to %v for %s: \n %v", p.Type, p.Value, p.Name, err)
}

// Set is a method implemented by Registry to configure the Windows registry
// according to the desired settings, while storing the previous values.
func (c Config) Set() error {
	for i, key := range c {
		k, err := registry.OpenKey(key.GetScope(), key.Path, registry.ALL_ACCESS)
		if err != nil {
			return err
		}
		for n, prop := range key.Properties {
			switch prop.Type {
			case "DWord":
				// First store the current val in prop.PrevVal,
				s, _, err := k.GetIntegerValue(prop.Name)
				if err != nil {
					return prop.getError(err)
				}
				c[i].Properties[n].PrevValue = s

				// Then set prop.Value to the current registry value.
				err = k.SetDWordValue(prop.Name, uint32(prop.Value.(float64)))
				if err != nil {
					return prop.setError(err)
				}
			case "QWord":
				// First store the current val in prop.PrevVal,
				s, _, err := k.GetIntegerValue(prop.Name)
				if err != nil {
					return prop.getError(err)
				}
				c[i].Properties[n].PrevValue = s

				// Then set prop.Value to the current registry value.
				err = k.SetQWordValue(prop.Name, uint64(prop.Value.(float64)))
				if err != nil {
					return prop.setError(err)
				}
			case "String":
				// First store the current val in prop.PrevVal,
				s, _, err := k.GetStringValue(prop.Name)
				if err != nil {
					return prop.getError(err)
				}
				c[i].Properties[n].PrevValue = s

				// Then set prop.Value to the current registry value.
				err = k.SetStringValue(prop.Name, prop.Value.(string))
				if err != nil {
					return prop.setError(err)
				}
			case "Strings":
				// First store the current val in prop.PrevVal,
				s, _, err := k.GetStringsValue(prop.Name)
				if err != nil {
					return prop.getError(err)
				}
				c[i].Properties[n].PrevValue = s

				// Then set prop.Value to the current registry value.
				err = k.SetStringsValue(prop.Name, prop.Value.([]string))
				if err != nil {
					return prop.setError(err)
				}
			case "Binary":
				// First store the current val in prop.PrevVal,
				s, _, err := k.GetBinaryValue(prop.Name)
				if err != nil {
					return prop.getError(err)
				}
				c[i].Properties[n].PrevValue = s

				// Then set prop.Value to the current registry value.
				err = k.SetBinaryValue(prop.Name, prop.Value.([]byte))
				if err != nil {
					return prop.setError(err)
				}
			default:
				return fmt.Errorf("please check the prop.Type is of the ValType type and try again")
			}
		}
		k.Close()
	}
	return nil
}
