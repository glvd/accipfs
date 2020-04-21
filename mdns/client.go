package mdns

import (
	"github.com/miekg/dns"
	"go.uber.org/atomic"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"strings"
	"time"
)

// ServiceEntry is returned after we query for a service
type ServiceEntry struct {
	Name       string
	Host       string
	AddrV4     net.IP
	AddrV6     net.IP
	Port       int
	Info       string
	InfoFields []string

	Addr net.IP // @Deprecated

	hasTXT bool
	sent   bool
}

// complete is used to check if we have all the info we need
func (s *ServiceEntry) complete() bool {
	return (s.AddrV4 != nil || s.AddrV6 != nil || s.Addr != nil) && s.Port != 0 && s.hasTXT
}

// QueryParam is used to customize how a Lookup is performed
type QueryParam struct {
	Service             string               // Service to lookup
	Domain              string               // Lookup domain, default "local"
	Timeout             time.Duration        // Lookup timeout, default 1 second
	Interface           *net.Interface       // Multicast interface to use
	Entries             chan<- *ServiceEntry // Entries Channel
	WantUnicastResponse bool                 // Unicast response desired, as per 5.4 in RFC
}

// DefaultParams is used to return a default set of QueryParam's
func DefaultParams(service string) *QueryParam {
	return &QueryParam{
		Service:             service,
		Domain:              "local",
		Timeout:             time.Second,
		Entries:             make(chan *ServiceEntry),
		WantUnicastResponse: false, // TODO(reddaly): Change this default.
	}
}

// Query looks up a given service, in a domain, waiting at most
// for a timeout before finishing the query. The results are streamed
// to a channel. Sends will not block, so clients should make sure to
// either read or buffer.
func (c *client) Query(params *QueryParam) error {
	// Set the multicast interface
	if params.Interface != nil {
		if err := c.setInterface(params.Interface); err != nil {
			return err
		}
	}

	// Ensure defaults are set
	if params.Domain == "" {
		params.Domain = "local"
	}
	if params.Timeout == 0 {
		params.Timeout = time.Second
	}

	// Run the query
	return c.query(params)
}

// Lookup is the same as Query, however it uses all the default parameters
func (c *client) Lookup(service string, entries chan<- *ServiceEntry) error {
	params := DefaultParams(service)
	params.Entries = entries
	return c.Query(params)
}

// Client ...
type Client interface {
	Query(params *QueryParam) error
	Lookup(service string, entries chan<- *ServiceEntry) error
}

type client struct {
	cfg     *OptionConfig
	uniConn []*net.UDPConn
	conn    []*net.UDPConn
	stop    *atomic.Bool
}

// Close ...
func (c *client) Close() (err error) {
	if !c.stop.CAS(false, true) {
		return
	}

	for i := range c.uniConn {
		if c.uniConn[i] != nil {
			e := c.uniConn[i].Close()
			if e != nil {
				err = e
			}
		}
	}
	for i := range c.conn {
		if c.conn[i] != nil {
			e := c.conn[i].Close()
			if e != nil {
				err = e
			}
		}
	}

	return err
}

// setInterface is used to set the query interface, uses system
// default if not provided
func (c *client) setInterface(iface *net.Interface) error {
	if c.uniConn[udp4] != nil {
		p := ipv4.NewPacketConn(c.uniConn[udp4])
		if err := p.SetMulticastInterface(iface); err != nil {
			return err
		}
	}
	if c.uniConn[udp6] != nil {
		p2 := ipv6.NewPacketConn(c.uniConn[udp6])
		if err := p2.SetMulticastInterface(iface); err != nil {
			return err
		}
	}
	if c.conn[udp4] != nil {
		p := ipv4.NewPacketConn(c.conn[udp4])
		if err := p.SetMulticastInterface(iface); err != nil {
			return err
		}
	}
	if c.conn[udp6] != nil {
		p2 := ipv6.NewPacketConn(c.conn[udp6])
		if err := p2.SetMulticastInterface(iface); err != nil {
			return err
		}
	}
	return nil
}

// query is used to perform a lookup and stream results
func (c *client) query(params *QueryParam) error {
	// Create the service name
	//serviceAddr := fmt.Sprintf("%s.%s.", trimDot(params.Service), trimDot(params.Domain))
	sa := serviceAddr(params.Service, params.Domain)
	// Start listening for response packets
	msgCh := make(chan *dns.Msg, 32)
	if c.uniConn[udp4] != nil {
		go c.recv(c.uniConn[udp4], msgCh)
	}
	if c.uniConn[udp6] != nil {
		go c.recv(c.uniConn[udp6], msgCh)
	}
	if c.conn[udp4] != nil {
		go c.recv(c.conn[udp4], msgCh)
	}

	if c.conn[udp6] != nil {
		go c.recv(c.conn[udp6], msgCh)
	}

	// Send the query
	m := new(dns.Msg)
	m.SetQuestion(sa, dns.TypePTR)
	// RFC 6762, section 18.12.  Repurposing of Top Bit of qclass in Question
	// Section
	//
	// In the Question Section of a Multicast DNS query, the top bit of the qclass
	// field is used to indicate that unicast responses are preferred for this
	// particular question.  (See Section 5.4.)
	if params.WantUnicastResponse {
		m.Question[0].Qclass |= 1 << 15
	}
	m.RecursionDesired = false
	if err := c.sendQuery(m); err != nil {
		return err
	}

	// Map the in-progress responses
	inprogress := make(map[string]*ServiceEntry)

	// Listen until we reach the timeout
	finish := time.After(params.Timeout)
	for {
		select {
		case resp := <-msgCh:
			var inp *ServiceEntry
			for _, answer := range append(resp.Answer, resp.Extra...) {
				// TODO(reddaly): Check that response corresponds to serviceAddr?
				switch rr := answer.(type) {
				case *dns.PTR:
					// Create new entry for this
					inp = ensureName(inprogress, rr.Ptr)

				case *dns.SRV:
					// Check for a target mismatch
					if rr.Target != rr.Hdr.Name {
						alias(inprogress, rr.Hdr.Name, rr.Target)
					}

					// Get the port
					inp = ensureName(inprogress, rr.Hdr.Name)
					inp.Host = rr.Target
					inp.Port = int(rr.Port)

				case *dns.TXT:
					// Pull out the txt
					inp = ensureName(inprogress, rr.Hdr.Name)
					inp.Info = strings.Join(rr.Txt, "|")
					inp.InfoFields = rr.Txt
					inp.hasTXT = true

				case *dns.A:
					// Pull out the IP
					inp = ensureName(inprogress, rr.Hdr.Name)
					inp.Addr = rr.A // @Deprecated
					inp.AddrV4 = rr.A

				case *dns.AAAA:
					// Pull out the IP
					inp = ensureName(inprogress, rr.Hdr.Name)
					inp.Addr = rr.AAAA // @Deprecated
					inp.AddrV6 = rr.AAAA
				}
			}

			if inp == nil {
				continue
			}

			// Check if this entry is complete
			if inp.complete() {
				if inp.sent {
					continue
				}
				inp.sent = true
				select {
				case params.Entries <- inp:
				default:
				}
			} else {
				// Fire off a node specific query
				m := new(dns.Msg)
				m.SetQuestion(inp.Name, dns.TypePTR)
				m.RecursionDesired = false
				if err := c.sendQuery(m); err != nil {
					logE("failed to query instance", "name", inp.Name, "error", err)
				}
			}
		case <-finish:
			return nil
		}
	}
}

// sendQuery is used to multicast a query out
func (c *client) sendQuery(q *dns.Msg) error {
	buf, err := q.Pack()
	if err != nil {
		return err
	}
	if c.uniConn[udp4] != nil {
		_, err = c.uniConn[udp4].WriteToUDP(buf, c.cfg.IPV4Addr)
		if err != nil {
			return err
		}
	}
	if c.uniConn[udp6] != nil {
		_, err = c.uniConn[udp6].WriteToUDP(buf, c.cfg.IPV6Addr)
		if err != nil {
			return err
		}
	}
	return nil
}

// recv is used to receive until we get a shutdown
func (c *client) recv(l *net.UDPConn, msgCh chan<- *dns.Msg) {
	if l == nil {
		return
	}
	buf := make([]byte, 65536)
	for !c.stop.Load() {
		n, err := l.Read(buf)

		if c.stop.Load() {
			return
		}

		if err != nil {
			logE("failed to read packet", "error", err)
			continue
		}
		msg := new(dns.Msg)
		if err := msg.Unpack(buf[:n]); err != nil {
			logE("failed to unpack packet", "error", err)
			continue
		}
		select {
		case msgCh <- msg:
		}
	}
}

// ensureName is used to ensure the named node is in progress
func ensureName(inprogress map[string]*ServiceEntry, name string) *ServiceEntry {
	if inp, ok := inprogress[name]; ok {
		return inp
	}
	inp := &ServiceEntry{
		Name: name,
	}
	inprogress[name] = inp
	return inp
}

// alias is used to setup an alias between two entries
func alias(inprogress map[string]*ServiceEntry, src, dst string) {
	srcEntry := ensureName(inprogress, src)
	inprogress[dst] = srcEntry
}
