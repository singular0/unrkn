package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
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
	flag.StringVar(&whitelistFilename, "w", "whitelist", "Whitelist file name. Expected format is one domain name per line. Domain will be converted to wildcard.")
}

func main() {
	flag.Parse()
	if outputFormat != "raw" && outputFormat != "routeros" {
		log.Fatal("Unknown format: ", outputFormat)
	}
	resp, err := http.Get(allURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	whitelist, err := whitelist.Load(whitelistFilename)
	if err != nil {
		log.Fatal(err)
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
			log.Fatalln(err)
		}
		if len(record) < 4 {
			continue
		}
		host := strings.TrimSpace(record[2])
		if whitelist.Contains(host) {
			ips := strings.Split(record[3], ",")
			for _, s := range ips {
				subnet := subnet.Parse(s)
				if subnet != nil {
					subnets.Add(*subnet)
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
			log.Fatalln(err)
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
