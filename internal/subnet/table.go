package subnet

import (
    "bufio"
    "os"
)

// IP4SubnetTable allows to sort IP4Subnet array
type IP4SubnetTable []IP4Subnet

func (a IP4SubnetTable) Len() int {
    return len(a)
}

func (a IP4SubnetTable) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

func (a IP4SubnetTable) Less(i, j int) bool {
    if a[i].MaskLength < a[j].MaskLength {
        return true
    }
    if a[i].MaskLength == a[j].MaskLength && a[i].Address < a[j].Address {
        return true
    }
    return false
}

// Add add subnet to the and merge it if contained by already existing subnet or contains some of them
func (a *IP4SubnetTable) Add(net IP4Subnet) {
    for i, n := range *a {
        if n.Contains(net) {
            return
        } else if net.Contains(n) {
            *a = append((*a)[:i], (*a)[i+1:]...)
        }
    }
    *a = append(*a, net)
}

// LoadTable loads ip table from file
func LoadTable(filename string) (*IP4SubnetTable, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    var table IP4SubnetTable
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        subnet := Parse(scanner.Text())
        if subnet != nil {
            table.Add(*subnet)
        }
    }
    return &table, scanner.Err()
}

// Save writes table contents to file
func (a IP4SubnetTable) Save(file *os.File) error {
    for _, subnet := range a {
        _, err := file.WriteString(subnet.String() + "\n")
        if err != nil {
            return err
        }
    }
    return nil
}

// ExportRouterOS export in RouterOS ip/firewall/address-list format
func (a IP4SubnetTable) ExportRouterOS(file *os.File, list string) error {
    for _, subnet := range a {
        s := "/ip firewall address-list add address=" + subnet.String()
        if (list != "") {
            s += " list=" + list
        }
        s += "\n"
        _, err := file.WriteString(s)
        if err != nil {
            return err
        }
    }
    return nil
}
