package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"net/http"
	"sort"
	"strings"

	"../../internal/subnet"
	"../../internal/whitelist"
)

const allURL = "http://api.antizapret.info/all.php"

var whitelistFilename string
var outputFilename string
var outputFormat string
var addressList string

func init() {
	flag.StringVar(&addressList, "a", "", "Address list name for routeros output format.")
	flag.StringVar(&outputFormat, "f", "raw", "Output format. One of: raw (one subnet per line), routeros (/ip/firewall/address-list add).")
	flag.StringVar(&outputFilename, "o", "", "Output file name. If omitted list will be sent to the console.")
	flag.StringVar(&whitelistFilename, "w", "", "Whitelist file name. Expected format is one domain name per line. Domain will be converted to wildcard.")
}

func main() {
	flag.Parse()
	if outputFormat != "raw" && outputFormat != "routeros" {
		msg := fmt.Sprintf("Unsupported output format: %s\n", outputFormat)
		os.Stderr.WriteString(msg)
		os.Exit(1)
	}
	resp, err := http.Get(allURL)
	if err != nil {
		msg := fmt.Sprintf("Error retrieving %s: %s\n", allURL, err)
		os.Stderr.WriteString(msg)
		os.Exit(1)
	}
	defer resp.Body.Close()
	var whiteList whitelist.Whitelist
	if whitelistFilename != "" {
		whiteList, err = whitelist.Load(whitelistFilename)
		if err != nil {
			msg := fmt.Sprintf("Error opening %s: %s\n", whitelistFilename, err)
			os.Stderr.WriteString(msg)
			os.Exit(1)
		}
	}
	r := csv.NewReader(resp.Body)
	r.Comma = ';'
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	r.ReuseRecord = true
	var subnets subnet.IP4SubnetTable
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			msg := fmt.Sprintf("Error reading data from %s: %s\n", allURL, err)
			os.Stderr.WriteString(msg)
			os.Exit(1)
		}
		if len(record) < 4 {
			continue
		}
		host := strings.TrimSpace(record[2])
		if host == "" || (whiteList != nil && whiteList.Contains(host)) {
			ips := strings.Split(record[3], ",")
			for _, s := range ips {
				subnet := subnet.Parse(s)
				if subnet != nil {
					subnets.Add(*subnet)
				} else {
					msg := fmt.Sprintf("Error parsing IP address %s\n", s)
					os.Stderr.WriteString(msg)
				}
			}
		}
	}
	sort.Sort(subnets)
	var file *os.File
	if outputFilename == "" {
		file = os.Stdout
	} else {
		file, err = os.Create(outputFilename)
		if err != nil {
			msg := fmt.Sprintf("Error writing data: %s\n", err)
			os.Stderr.WriteString(msg)
			os.Exit(1)
		}
		defer file.Close()
	}
	switch(outputFormat) {
	case "routeros":
		subnets.ExportRouterOS(file, addressList)
	case "raw":
		subnets.Save(file)
	}
}
