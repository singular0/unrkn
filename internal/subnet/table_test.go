package subnet

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct{
		s string
		n int
	}{
		{ "192.168.0.1", 1 },
		{ "192.168.0.1", 1 },
		{ "192.168.0.2", 2 },
		{ "192.168.0.2", 2 },
		{ "192.168.0.3", 3 },
		{ "192.168.0.2", 3 },
		{ "192.168.0.1", 3 },
	}

	table := IP4SubnetTable{}
	for i, test := range tests {
		subnet := Parse(test.s)
		table.Add(*subnet)
		len := table.Len()
		if len != test.n {
			t.Errorf("Failed Add() for entry %d ('%s'): array length is %d, should be %d", i, test.s, len, test.n)
		}
	}
}
