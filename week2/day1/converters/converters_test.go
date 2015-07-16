package converters

import "testing"

func TestConvert(t *testing.T) {
	cases := []struct{ from, to, expected string }{
		{"50mi", "km", "80.47km"},
		{"50km", "mi", "31.07mi"},
	}

	for _, c := range cases {
		str, err := Convert(c.from, c.to)
		if err != nil {
			t.Log("error should be nil", err)
			t.Fail()
		}
		if str != c.expected {
			t.Log("error should be "+c.expected+" but got", str)
			t.Fail()
		}
	}

}
