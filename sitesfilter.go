package filterschedule

import (
	"strings"
	"time"

	"os"

	"k8s.io/apimachinery/pkg/util/yaml"
)

type SitesFilter struct {
	Description string
	Filter      Sites
	For         Targets
	When        Schedule
}

type Sites struct {
	Everything bool
	Sites      []string
}

type Targets struct {
	All   bool
	Hosts []string
}

func LoadFromFile(filename string) ([]SitesFilter, error) {
	var f []SitesFilter
	content, err := os.ReadFile(filename)
	if err == nil {
		err = yaml.Unmarshal(content, &f)
	}
	return f, err
}

func (sf *SitesFilter) IsMatching(name string, clientIP string, t time.Time) bool {
	if sf.Filter.Everything {
		if sf.isApplicableTo(clientIP) && sf.isApplicableAt(t) {
			return true
		}
	} else {
		for _, p := range sf.Filter.Sites {
			if strings.Contains(name, p) {
				if sf.isApplicableTo(clientIP) && sf.isApplicableAt(t) {
					return true
				}
			}
		}
	}
	return false
}

func (sf *SitesFilter) isApplicableTo(clientIP string) bool {
	if sf.For.All {
		return true
	} else {
		for _, ip := range sf.For.Hosts {
			if ip == clientIP {
				return true
			}
		}
	}
	return false
}

func (sf *SitesFilter) isApplicableAt(t time.Time) bool {
	if sf.When.Always {
		return true
	} else {
		for _, se := range sf.When.Schedule {
			if isApplicableAtDays(se.Days, t) && isApplicableAtHours(se.Hours, t) {
				return true
			}
		}
	}
	return false
}
