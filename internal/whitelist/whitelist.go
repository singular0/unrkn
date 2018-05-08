package whitelist

import (
	"bufio"
	"os"
	"regexp"
)

// Whitelist contains list of whitelisted domains
type Whitelist []*regexp.Regexp

// LoadWhitelist loads whitelisted domains from the text file
func Load(filename string) (Whitelist, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	whitelist := Whitelist{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := ".*\\." + scanner.Text() + "$"
		re := regexp.MustCompile(s)
		if re != nil {
			whitelist = append(whitelist, re)
		}
	}
	if err := scanner.Err(); err != nil {
		return whitelist, err
	}
	return whitelist, nil
}

// Contains check if domain contained in the whitelist
func (w Whitelist) Contains(host string) bool {
	if host == "" {
		return true
	}
	for _, re := range w {
		if re.MatchString(host) {
			return true
		}
	}
	return false
}
