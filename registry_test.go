package regselect

import (
	"testing"

	"golang.org/x/sys/windows/registry"
)

func TestGetScope(t *testing.T) {
	cases := []struct {
		in   string
		want registry.Key
	}{
		{"LOCAL_MACHINE", registry.LOCAL_MACHINE},
		{"CURRENT_USER", registry.CURRENT_USER},
		{"CLASSES_ROOT", registry.CLASSES_ROOT},
		{"CURRENT_CONFIG", registry.CURRENT_CONFIG},
		{"USERS", registry.USERS},
		{"", registry.LOCAL_MACHINE},
	}
	for _, c := range cases {
		k := Key{}
		k.Scope = c.in
		got := k.GetScope()
		if got != c.want {
			t.Errorf("Method GetScope() implemented incorrectly on key: in:  %q \n got: %q \n want: %q \n", c.in, got, c.want)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	var c Config
	fi := "config_test.json"
	c, err := Unmarshal(fi)
	if err != nil {
		t.Errorf("unable to unmarshal file %s: %#v", fi, err)
	}
	if l := len(c); l != 1 {
		t.Errorf("Config has wrong number of keys: %d (length should be 1)")
	}
	for _, v := range c {
		n := 0
		for _, k := range v {
			switch k.Name {
			case "ProxyEnable":
				if k.Type != "DWord" || k.Value != 0 {
					t.Errorf("unmarshaled incorrectly: %s returned %s type and %d value! \n", k.Name, k.Type, k.Value)
				}
				n++
			case "ProxyServer":
				if k.Type != "String" || k.Value != "" {
					t.Errorf("unmarshaled incorrectly: %s returned %s type and %s value! \n", k.Name, k.Type, k.Value)
				}
				n++
			case "WarnOnIntranet":
				if k.Type != "DWord" || k.Value != 0 {
					t.Errorf("unmarshaled incorrectly: %s returned %s type and %d value! \n", k.Name, k.Type, k.Value)
				}
				n++
			}
		}
		if n != 3 {
			t.Errorf("failed to unmarshal all of the properties listed in the .json test file: got %d expected %d", n, 3)
		}
	}
}

func TestWrite(t *testing.T) {
	var c Config
	fi := "config_test.json"
	c, err := Unmarshal(fi)
	if err != nil {
		t.Errorf("unable to unmarshal file %s", fi)
	}

	if err := c.Write("TestWrite_test.json"); err != nil {
		t.Errorf("couldn't write config: %#v", err)
	}

	fs := []string{fi, "TestWrite_test.json"}
	cs := make([]Config, 1)
	for _, v := range fs {
		b, err := Unmarshal(v)
		if err != nil {
			t.Errorf("couldn't unmarshal test JSON file %s: %#v \n", v, err)
		}
		cs = append(cs, b)
	}
	if err := len(cs); err != 1 {
		t.Errorf("not enough Config structs were unmarshaled: need 2")
	}

	if cs[0] != cs[1] {
		t.Errorf("config %#v is not equal to config %#v", cs[0], cs[1])
	}
}
