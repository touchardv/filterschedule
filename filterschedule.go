package filterschedule

import (
	"context"
	"time"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("filterschedule")

type FilterSchedule struct {
	filters []SitesFilter
	Next    plugin.Handler
}

// ServeDNS implements the plugin.Handler interface. This method gets called when filterschedule is used
// in a Server.
func (e FilterSchedule) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	req := request.Request{W: w, Req: r}
	timestamp := time.Now()
	for _, f := range e.filters {
		if f.IsMatching(req.Name(), req.IP(), timestamp) {
			log.Warningf("Filter '%s' matched: %s", f.Description, req.Name())
			m := newResponseMessage(req)
			m.SetReply(r)
			m.SetRcode(r, dns.RcodeSuccess)
			w.WriteMsg(m)
			return dns.RcodeSuccess, nil
		}
	}

	// Call next plugin (if any).
	return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (e FilterSchedule) Name() string { return "filterschedule" }
