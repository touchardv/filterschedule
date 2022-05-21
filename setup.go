package filterschedule

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("filterschedule", setup) }

func setup(c *caddy.Controller) error {
	cfgFilename, err := parse(c)
	if err != nil {
		return plugin.Error("filterschedule", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		f, err := LoadFromFile(cfgFilename)
		if err != nil {
			log.Warning("Failed to load: ", err)
			f = make([]SitesFilter, 0)
		} else {
			log.Infof("Loaded %d filter(s)", len(f))
		}
		return FilterSchedule{
			filters: f,
			Next:    next,
		}
	})

	// All OK, return a nil error.
	return nil
}

func parse(c *caddy.Controller) (string, error) {
	const defaultCfgFilename = "filterschedule.yaml"
	for c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 0:
			return defaultCfgFilename, nil
		case 1:
			return args[0], nil
		default:
			return "", plugin.Error("filterschedule", c.ArgErr())
		}
	}
	return "", plugin.Error("filterschedule", c.ArgErr())
}
