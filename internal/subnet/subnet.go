package subnet

import (
    "fmt"
    "regexp"
    "strconv"
)

// IP4Subnet is a representation of IPv4 subnet
type IP4Subnet struct {
    Address    uint32
    Mask       uint32
    MaskLength uint8
}

var parser = regexp.MustCompile(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})(?:/(\d{1,2}))?$`)

// Parse parse any CIDR subnet or single IP address from string
func Parse(s string) *IP4Subnet {
    parts := parser.FindStringSubmatch(s)
    if len(parts) == 0 {
        return nil
    }
    subnet := IP4Subnet{}
    if parts[5] != "" {
        bits, _ := strconv.ParseUint(parts[5], 10, 8)
        subnet.MaskLength = uint8(bits)
    } else {
        subnet.MaskLength = 32
    }
    subnet.Mask = 0xFFFFFFFF << (32 - subnet.MaskLength)
    for i := 1; i < 5; i++ {
        byte, _ := strconv.ParseUint(parts[i], 10, 8)
        subnet.Address |= uint32(byte) << uint8((4 - i) * 8)
    }
    return &subnet
}

// Contains check if subnet contains or equal to other subnet
func (n IP4Subnet) Contains(subnet IP4Subnet) bool {
    return n.MaskLength <= subnet.MaskLength && (n.Address & n.Mask) == (subnet.Address & n.Mask)
}

// IsPrivate checks if subnet contained or equal to RFC1918 private nets
func (n IP4Subnet) IsPrivate() bool {
    // 192.168.0.0/16
    if n.Address & 0xFFFF0000 == 0xC0A80000 {
        return true
    }
    // 172.16.0.0/12
    if n.Address & 0xFFF00000 == 0xAC100000 {
        return true
    }
    // 10.0.0.0/8
    if n.Address & 0xFF000000 == 0x0A000000 {
        return true
    }
    return false
}

// IsLocal checks if subnet contained or equal to RFC1122 & RFC3927 defined nets
func (n IP4Subnet) IsLocal() bool {
    // 0.0.0.0/8
    if n.Address & 0xFF000000 == 0x00000000 {
        return true
    }
    // 127.0.0.0/8
    if n.Address & 0xFF000000 == 0x7F000000 {
        return true
    }
    // 169.254.0.0/16
    if n.Address & 0xFFFF0000 == 0xA9FE0000 {
        return true
    }
    return false
}

func (n IP4Subnet) String() string {
    s := fmt.Sprintf("%d.%d.%d.%d", uint8(n.Address >> 24), uint8(n.Address >> 16), uint8(n.Address >> 8), uint8(n.Address))
    if n.MaskLength != 32 {
        s += fmt.Sprintf("/%d", n.MaskLength)
    }
    return s
}
