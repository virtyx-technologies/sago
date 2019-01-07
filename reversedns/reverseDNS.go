package reversedns
// TODO : deprecated, see net.Resolver
import (
	"bufio"
	"github.com/virtyx-technologies/sago/globals"
	"github.com/virtyx-technologies/sago/util"
	"log"
	"os/exec"
	"strings"
)

type (
	filter interface {
		filter(string) (bool, string)
	}

	resolver struct {
		command string   //  Command to execute; e.g. dig
		args    []string // Arguments for the command
		ipIndex int      // Index of argument representing ip address
		filter  string   // Name of the filter used to extract the result
	}
)

var (
	localResolver *resolver
	stdResolver   *resolver
)

// Return the host name for the given IP
func ReverseDNS(ip string) string {
	var host = "Unknown"
	var r *resolver
	if strings.HasSuffix(ip, ".local") {
		r = localResolver
	} else {
		r = stdResolver
	}

	r.args[r.ipIndex] = ip
	cmdOut, err := exec.Command(r.command, r.args...).Output()
	if err != nil {
		log.Printf("There was an error running %s command: %s", r.command, err)
		return host
	}

	filter := filters[r.filter].(filter)
	var match bool
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOut)))
	for scanner.Scan() {
		match, host = filter.filter(scanner.Text())
		if match {
			break
		}
	}
	return host
}

// Establish local and standard resolvers
func init() {

	if found, _ := util.IsOnPath("avahi-resolve"); found {
		localResolver = NewResolver("avahi-resolve -a %s", "filterAvahi")
	} else if found, _ := util.IsOnPath("dig"); found {
		localResolver = NewResolver("dig -x %s @224.0.0.251 -p 5353 +notcp +noall +answer", "filterDig")
	}

	if found, _ := util.IsOnPath("dig"); found {
		stdResolver = NewResolver("dig -x %s +noall +answer", "filterDig")
	} else if found, _ := util.IsOnPath("host"); found {
		stdResolver = NewResolver("host -t PTR %s", "filterHost")
	} else if found, _ := util.IsOnPath("drill"); found {
		stdResolver = NewResolver("drill -x ptr %s", "filterDrill")
	} else if found, _ := util.IsOnPath("nslookup"); found {
		stdResolver = NewResolver(`nslookup -type=PTR %s`, "filterNslookup")
	}

	if localResolver == nil {
		log.Fatal("Neither \"dig\"  nor \"avahi-resolve\" is present", globals.ERR_DNSBIN)
	}

	if stdResolver == nil {
		log.Fatal("Neither \"dig\", \"host\", \"drill\" nor \"nslookup\" is present", globals.ERR_DNSBIN)
	}

}

func NewResolver(cmdline string, filter string) *resolver {
	tokens := strings.Split(cmdline, " ")
	command := tokens[0]
	var args []string
	copy(args, tokens[1:])
	var index int
	for i, arg := range args {
		if arg == "%s" {
			index = i
			break
		}
	}

	return &resolver{command, args, index, filter}
}
