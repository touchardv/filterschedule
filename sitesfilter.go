package filterschedule

import (
	"io/ioutil"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/yaml"
)

type SitesFilter struct {
	Description string
	Sites       []string
	For         Targets
	When        Schedule
}

type Targets struct {
	All   bool
	Hosts []string
}

func LoadFromFile(filename string) ([]SitesFilter, error) {
	var f []SitesFilter
	content, err := ioutil.ReadFile(filename)
	if err == nil {
		err = yaml.Unmarshal(content, &f)
	}
	return f, err
}

func (sf *SitesFilter) IsMatching(name string, clientIP string, t time.Time) bool {
	for _, p := range sf.Sites {
		if strings.Contains(name, p) {
			if sf.isApplicableTo(clientIP) && sf.isApplicableAt(t) {
				return true
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
