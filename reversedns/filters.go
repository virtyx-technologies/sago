package reversedns

import "strings"

// Filter instances indexed by name
var filters map[string]interface{} = map[string]interface{}{
	"filterAvahi":    filterAvahi{},
	"filterNslookup": filterNslookup{},
	"filterDig":      filterDig{},
	"filterHost":     filterHost{},
	"filterDrill":    filterDrill{},
}

// Filter for 'avahi-resolve' output
type filterAvahi struct{}

func (x *filterAvahi) filter(input string) (bool, string) {
	// Original filter:  awk '{ print $2 }'"
	fields := strings.Fields(input)
	if len(fields) > 1 {
		return true, fields[1]
	} else {
		return false, ""
	}
}

// Filter for 'dig' output
type filterDig struct{}

func (x *filterDig) filter(input string) (bool, string) {
	// Original filter:  awk  '/PTR/ { print $NF }'
	if strings.Contains(input, "PTR") {
		fields := strings.Fields(input)
		return true, fields[len(fields)-1]
	} else {
		return false, ""
	}
}

// Filter for 'host' output
type filterHost struct{}

func (x *filterHost) filter(input string) (bool, string) {
	// Original filter:  awk '/pointer/ { print $NF }'
	if strings.Contains(input, "pointer") {
		fields := strings.Fields(input)
		return true, fields[len(fields)-1]
	} else {
		return false, ""
	}
}

// Filter for 'drill' output
type filterDrill struct{}

func (x *filterDrill) filter(input string) (bool, string) {
	// Original filter:  awk '/ANSWER SECTION/ { getline; print $NF }'
}

// Filter for 'nslookup' output
type filterNslookup struct{}

func (x *filterNslookup) filter(input string) (bool, string) {
	// Original filter:  grep -v 'canonical name =' grep 'name = ' awk '{ print $NF }' sed 's/\.$//')
}
