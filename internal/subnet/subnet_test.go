package subnet

import "testing"

func TestParse(t *testing.T) {
    tests := []struct{
        s string
        result *IP4Subnet
    }{
        { "192.168.0.0",	&IP4Subnet{ Address: 0xC0A80000, Mask: 0xFFFFFFFF, MaskLength: 32 } },
        { "192.168.0.0/24", &IP4Subnet{ Address: 0xC0A80000, Mask: 0xFFFFFF00, MaskLength: 24 } },
        { "foobar",			nil },
    }

    for _, test := range tests {
        subnet := Parse(test.s)
        if subnet != nil {
            if test.result == nil {
                t.Errorf("Failed Parse() for '%s': should be nil", test.s)
            } else if *subnet != *test.result {
                t.Errorf("Failed Parse() for '%s': Address=%04X Mask=%04X MaskLength=%d", test.s, subnet.Address, subnet.Mask, subnet.MaskLength)
            }
        } else if subnet == nil && test.result != nil {
            t.Errorf("Failed Parse() for '%s': should not be nil", test.s)
        }
    }
}

func TestString(t *testing.T) {
    tests := []struct{
        subnet IP4Subnet
        result string 
    }{
        { IP4Subnet{ Address: 0xC0A80000, Mask: 0xFFFFFFFF, MaskLength: 32 }, "192.168.0.0" },
        { IP4Subnet{ Address: 0xC0A80000, Mask: 0xFFFFFF00, MaskLength: 24 }, "192.168.0.0/24" },
    }
    
    for _, test := range tests {
        s := test.subnet.String()
        if s != test.result {
            t.Errorf("Failed String() for '%s': should be '%s'", test.subnet, test.result)
        }
    }
}

func TestContains(t *testing.T) {
    tests := []struct{
        s1 string
        s2 string
        result bool
    }{
        { "192.168.0.0/24", "192.168.0.1", 		true },
        { "192.168.0.0/24", "192.167.0.1", 		false },
        { "192.0.0.0/8", 	"192.168.0.0/24", 	true },
        { "192.168.0.1", 	"192.168.0.1", 		true },
        { "192.168.0.1", 	"192.168.0.1/32",	true },
        { "192.168.0.3", 	"192.168.0.1",		false },
    }
    
    for _, test := range tests {
        subnet1 := Parse(test.s1)
        subnet2 := Parse(test.s2)
        if subnet1.Contains(*subnet2) != test.result {
            t.Errorf("Failed %s.Contains(%s): should be %t", subnet1, subnet2, test.result)
        }
    }
}

func TestIsPrivate(t *testing.T) {
    tests := []struct{
        s string
        result bool
    }{
        { "192.168.0.0/24", true },
        { "192.168.0.0/8",	true },
        { "10.10.10.1",		true },
        { "8.8.8.8",		false },
        { "8.8.8.8/8",		false },
        { "192.167.0.0/16",	false },
    }
    
    for _, test := range tests {
        subnet := Parse(test.s)
        if subnet.IsPrivate() != test.result {
            t.Errorf("Failed %s.IsPrivate(): should be %t", subnet, test.result)
        }
    }
}
