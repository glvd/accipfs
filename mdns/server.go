package mdns

import (
	"fmt"
	"github.com/miekg/dns"
	"go.uber.org/atomic"
	"math/rand"
	"net"
	"strings"
	"time"
)

const forceUnicastResponses = false

// Server ...
type Server interface {
	Start()
	Stop() error
}

type server struct {
	serviceAddr  string
	instanceAddr string
	enumAddr     string
	cfg          *OptionConfig
	conn         []*net.UDPConn
	stop         *atomic.Bool
	//wg           sync.WaitGroup
}

// Stop ...
func (s *server) Stop() (err error) {
	if !s.stop.CAS(false, true) {
		return nil
	}
	if s.conn[udp4] != nil {
		e := s.conn[udp4].Close()
		if e != nil {
			err = e
		}
	}
	if s.conn[udp6] != nil {
		e := s.conn[udp6].Close()
		if e != nil {
			err = e
		}
	}
	return err
}

// Start ...
func (s *server) Start() {
	for i := range s.conn {
		if s.conn[i] != nil {
			go s.recv(s.conn[i])
		}
	}
	go s.probe()
}

// recv is a long running routine to receive packets from an interface
func (s *server) recv(c *net.UDPConn) {
	buf := make([]byte, 65536)
	for !s.stop.Load() {
		//logI("reading from remote conn")
		n, from, err := c.ReadFrom(buf)
		if err != nil {
			//logE("failed to read from buffer")
			continue
		}
		//logI("parse from", "addr", from.String())
		if err := s.parsePacket(buf[:n], from); err != nil {
			//logE("failed to handle query", "error", err)
		}
	}
}

// parsePacket is used to parse an incoming packet
func (s *server) parsePacket(packet []byte, from net.Addr) error {
	var msg dns.Msg
	if err := msg.Unpack(packet); err != nil {
		//logE("failed to unpack packet", "error", err)
		return err
	}
	return s.handleQuery(&msg, from)
}

// handleQuery is used to handle an incoming query
func (s *server) handleQuery(query *dns.Msg, from net.Addr) error {
	if query.Opcode != dns.OpcodeQuery {
		// "In both multicast query and multicast response messages, the OPCODE MUST
		// be zero on transmission (only standard queries are currently supported
		// over multicast).  Multicast DNS messages received with an OPCODE other
		// than zero MUST be silently ignored."  Note: OpcodeQuery == 0
		return fmt.Errorf("mdns: received query with non-zero Opcode %v: %v", query.Opcode, *query)
	}
	if query.Rcode != 0 {
		// "In both multicast query and multicast response messages, the Response
		// Code MUST be zero on transmission.  Multicast DNS messages received with
		// non-zero Response Codes MUST be silently ignored."
		return fmt.Errorf("mdns: received query with non-zero Rcode %v: %v", query.Rcode, *query)
	}

	// TODO(reddaly): Handle "TC (Truncated) Bit":
	//    In query messages, if the TC bit is set, it means that additional
	//    Known-Answer records may be following shortly.  A responder SHOULD
	//    record this fact, and wait for those additional Known-Answer records,
	//    before deciding whether to respond.  If the TC bit is clear, it means
	//    that the querying host has no additional Known Answers.
	if query.Truncated {
		return fmt.Errorf("[ERR] mdns: support for DNS requests with high truncated bit not implemented: %v", *query)
	}

	var unicastAnswer, multicastAnswer []dns.RR

	// Handle each question
	for _, q := range query.Question {
		mrecs, urecs := s.handleQuestion(q)
		multicastAnswer = append(multicastAnswer, mrecs...)
		unicastAnswer = append(unicastAnswer, urecs...)
	}

	// See section 18 of RFC 6762 for rules about DNS headers.
	resp := func(unicast bool) *dns.Msg {
		// 18.1: ID (Query Identifier)
		// 0 for multicast response, query.Id for unicast response
		id := uint16(0)
		if unicast {
			id = query.Id
		}

		var answer []dns.RR
		if unicast {
			answer = unicastAnswer
		} else {
			answer = multicastAnswer
		}
		if len(answer) == 0 {
			return nil
		}

		return &dns.Msg{
			MsgHdr: dns.MsgHdr{
				Id: id,

				// 18.2: QR (Query/Response) Bit - must be set to 1 in response.
				Response: true,

				// 18.3: OPCODE - must be zero in response (OpcodeQuery == 0)
				Opcode: dns.OpcodeQuery,

				// 18.4: AA (Authoritative Answer) Bit - must be set to 1
				Authoritative: true,

				// The following fields must all be set to 0:
				// 18.5: TC (TRUNCATED) Bit
				// 18.6: RD (Recursion Desired) Bit
				// 18.7: RA (Recursion Available) Bit
				// 18.8: Z (Zero) Bit
				// 18.9: AD (Authentic Data) Bit
				// 18.10: CD (Checking Disabled) Bit
				// 18.11: RCODE (Response Code)
			},
			// 18.12 pertains to questions (handled by handleQuestion)
			// 18.13 pertains to resource records (handled by handleQuestion)

			// 18.14: Name Compression - responses should be compressed (though see
			// caveats in the RFC), so set the Compress bit (part of the dns library
			// API, not part of the DNS packet) to true.
			Compress: true,

			Answer: answer,
		}
	}
	//logI("query detail", "id", query.Id, "question", query.Question, "answer", query.Answer)
	if s.cfg.LogEmptyResponses && len(multicastAnswer) == 0 && len(unicastAnswer) == 0 {
		questions := make([]string, len(query.Question))
		for i, q := range query.Question {
			questions[i] = q.Name
		}
		logE("no responses for query with questions", "question", strings.Join(questions, ", "))
	}

	if mresp := resp(false); mresp != nil {
		//logI("multicast", "response", *mresp)
		if err := s.sendResponse(mresp, from, false); err != nil {
			return fmt.Errorf("mdns: error sending multicast response: %v", err)
		}
	}
	if uresp := resp(true); uresp != nil {
		if err := s.sendResponse(uresp, from, true); err != nil {
			//logI("unicast", "response", *uresp)
			return fmt.Errorf("mdns: error sending unicast response: %v", err)
		}
	}
	return nil
}

// handleQuestion is used to handle an incoming question
//
// The response to a question may be transmitted over multicast, unicast, or
// both.  The return values are DNS records for each transmission type.
func (s *server) handleQuestion(q dns.Question) (multicastRecs, unicastRecs []dns.RR) {
	records := s.Records(q)

	if len(records) == 0 {
		return nil, nil
	}

	// Handle unicast and multicast responses.
	// TODO(reddaly): The decision about sending over unicast vs. multicast is not
	// yet fully compliant with RFC 6762.  For example, the unicast bit should be
	// ignored if the records in question are close to TTL expiration.  For now,
	// we just use the unicast bit to make the decision, as per the spec:
	//     RFC 6762, section 18.12.  Repurposing of Top Bit of qclass in Question
	//     Section
	//
	//     In the Question Section of a Multicast DNS query, the top bit of the
	//     qclass field is used to indicate that unicast responses are preferred
	//     for this particular question.  (See Section 5.4.)
	if q.Qclass&(1<<15) != 0 || forceUnicastResponses {
		return nil, records
	}
	return records, nil
}

// sendResponse is used to send a response packet
func (s *server) sendResponse(resp *dns.Msg, from net.Addr, unicast bool) error {
	// TODO(reddaly): Respect the unicast argument, and allow sending responses
	// over multicast.
	buf, err := resp.Pack()
	if err != nil {
		return err
	}

	// Determine the socket to send from
	addr := from.(*net.UDPAddr)
	if addr.IP.To4() != nil {
		_, err = s.conn[udp4].WriteToUDP(buf, addr)
		return err
	}
	_, err = s.conn[udp6].WriteToUDP(buf, addr)
	return err
}

// Records ...
func (s *server) Records(q dns.Question) []dns.RR {
	list := map[string]func(question dns.Question) []dns.RR{
		s.cfg.enumAddr:     s.enumRecords,
		s.cfg.serviceAddr:  s.serviceRecords,
		s.cfg.instanceAddr: s.instanceRecords,
		s.cfg.HostName:     s.instanceRecords,
	}

	f, b := list[q.Name]
	if b {
		if q.Name == s.cfg.HostName {
			if q.Qtype == dns.TypeA || q.Qtype == dns.TypeAAAA {
				//do nothing
			} else {
				return nil
			}
		}
		//else if q.Name == "_services._dns-sd._udp."+s.cfg.Domain+"." {
		//	recs = s.dnssdMetaQueryRecords(q)
		//}
		//if recs != nil {
		//	return append(recs, f(q)...)
		//}
		return f(q)
	}
	return nil
}

func (s *server) enumRecords(q dns.Question) []dns.RR {
	switch q.Qtype {
	case dns.TypeANY:
		fallthrough
	case dns.TypePTR:
		rr := &dns.PTR{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    s.cfg.TTL,
			},
			Ptr: s.cfg.serviceAddr,
		}
		return []dns.RR{rr}
	default:
		return nil
	}
}

func (s *server) serviceRecords(q dns.Question) []dns.RR {
	switch q.Qtype {
	case dns.TypeANY:
		fallthrough
	case dns.TypePTR:
		// Build a PTR response for the service
		rr := &dns.PTR{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    s.cfg.TTL,
			},
			Ptr: s.cfg.instanceAddr,
		}
		servRec := []dns.RR{rr}

		// Get the instance records
		instRecs := s.instanceRecords(dns.Question{
			Name:  s.cfg.instanceAddr,
			Qtype: dns.TypeANY,
		})

		// Return the service record with the instance records
		return append(servRec, instRecs...)
	default:
		return nil
	}
}

func (s *server) instanceRecords(q dns.Question) []dns.RR {
	switch q.Qtype {
	case dns.TypeANY:
		// Get the SRV, which includes A and AAAA
		recs := s.instanceRecords(dns.Question{
			Name:  s.cfg.instanceAddr,
			Qtype: dns.TypeSRV,
		})

		// Add the TXT record
		recs = append(recs, s.instanceRecords(dns.Question{
			Name:  s.cfg.instanceAddr,
			Qtype: dns.TypeTXT,
		})...)
		return recs

	case dns.TypeA:
		var rr []dns.RR
		for _, ip := range s.cfg.IPs {
			if ip4 := ip.To4(); ip4 != nil {
				rr = append(rr, &dns.A{
					Hdr: dns.RR_Header{
						Name:   s.cfg.HostName,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    s.cfg.TTL,
					},
					A: ip4,
				})
			}
		}
		return rr

	case dns.TypeAAAA:
		var rr []dns.RR
		for _, ip := range s.cfg.IPs {
			if ip.To4() != nil {
				// TODO(reddaly): IPv4 addresses could be encoded in IPv6 format and
				// putinto AAAA records, but the current logic puts ipv4-encodable
				// addresses into the A records exclusively.  Perhaps this should be
				// configurable?
				continue
			}

			if ip16 := ip.To16(); ip16 != nil {
				rr = append(rr, &dns.AAAA{
					Hdr: dns.RR_Header{
						Name:   s.cfg.HostName,
						Rrtype: dns.TypeAAAA,
						Class:  dns.ClassINET,
						Ttl:    s.cfg.TTL,
					},
					AAAA: ip16,
				})
			}
		}
		return rr

	case dns.TypeSRV:
		// Create the SRV Record
		srv := &dns.SRV{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeSRV,
				Class:  dns.ClassINET,
				Ttl:    s.cfg.TTL,
			},
			Priority: 10,
			Weight:   1,
			Port:     s.cfg.Port,
			Target:   s.cfg.HostName,
		}
		recs := []dns.RR{srv}

		// Add the A record
		recs = append(recs, s.instanceRecords(dns.Question{
			Name:  s.cfg.instanceAddr,
			Qtype: dns.TypeA,
		})...)

		// Add the AAAA record
		recs = append(recs, s.instanceRecords(dns.Question{
			Name:  s.cfg.instanceAddr,
			Qtype: dns.TypeAAAA,
		})...)
		return recs

	case dns.TypeTXT:
		txt := &dns.TXT{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassINET,
				Ttl:    s.cfg.TTL,
			},
			Txt: s.cfg.TXT,
		}
		return []dns.RR{txt}
	}
	return nil
}

// dnssdMetaQueryRecords returns the DNS records in response to a "meta-query"
// issued to browse for DNS-SD services, as per section 9. of RFC6763.
//
// A meta-query has a name of the form "_services._dns-sd._udp.<Domain>" where
// Domain is a fully-qualified domain, such as "local."
func (s *server) dnssdMetaQueryRecords(q dns.Question) []dns.RR {
	// Intended behavior, as described in the RFC:
	//     ...it may be useful for network administrators to find the list of
	//     advertised service types on the network, even if those Service Names
	//     are just opaque identifiers and not particularly informative in
	//     isolation.
	//
	//     For this purpose, a special meta-query is defined.  A DNS query for PTR
	//     records with the name "_services._dns-sd._udp.<Domain>" yields a set of
	//     PTR records, where the rdata of each PTR record is the two-abel
	//     <Service> name, plus the same domain, e.g., "_http._tcp.<Domain>".
	//     Including the domain in the PTR rdata allows for slightly better name
	//     compression in Unicast DNS responses, but only the first two labels are
	//     relevant for the purposes of service type enumeration.  These two-label
	//     service types can then be used to construct subsequent Service Instance
	//     Enumeration PTR queries, in this <Domain> or others, to discover
	//     instances of that service type.
	return []dns.RR{
		&dns.PTR{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypePTR,
				Class:  dns.ClassINET,
				Ttl:    defaultTTL,
			},
			Ptr: s.serviceAddr,
		},
	}
}
func (s *server) probe() {
	//defer s.wg.Done()

	name := fmt.Sprintf("%s.%s.%s.", s.cfg.Instance, trimDot(s.cfg.Service), trimDot(s.cfg.Domain))

	q := new(dns.Msg)
	q.SetQuestion(name, dns.TypePTR)
	q.RecursionDesired = false

	srv := &dns.SRV{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeSRV,
			Class:  dns.ClassINET,
			Ttl:    defaultTTL,
		},
		Priority: 0,
		Weight:   0,
		Port:     s.cfg.Port,
		Target:   s.cfg.HostName,
	}
	txt := &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    defaultTTL,
		},
		Txt: s.cfg.TXT,
	}
	q.Ns = []dns.RR{srv, txt}

	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 3; i++ {
		if err := s.SendMulticast(q); err != nil {
			logI("failed to send probe", "error", err)
		}
		time.Sleep(time.Duration(randomizer.Intn(250)) * time.Millisecond)
	}

	resp := new(dns.Msg)
	resp.MsgHdr.Response = true

	// set for query
	q.SetQuestion(name, dns.TypeANY)

	resp.Answer = append(resp.Answer, s.Records(q.Question[0])...)

	// reset
	q.SetQuestion(name, dns.TypePTR)

	// From RFC6762
	//    The Multicast DNS responder MUST send at least two unsolicited
	//    responses, one second apart. To provide increased robustness against
	//    packet loss, a responder MAY send up to eight unsolicited responses,
	//    provided that the interval between unsolicited responses increases by
	//    at least a factor of two with every response sent.
	timeout := 1 * time.Second
	timer := time.NewTimer(timeout)
	for i := 0; i < 3; i++ {
		if err := s.SendMulticast(resp); err != nil {
			logE("failed to send announcement", "error", err)
		}
		select {
		case <-timer.C:
			timeout *= 2
			timer.Reset(timeout)
		default:
			if s.stop.Load() {
				return
			}
		}
	}
}

// multicastResponse us used to send a multicast response packet
func (s *server) SendMulticast(msg *dns.Msg) error {
	buf, err := msg.Pack()
	if err != nil {
		return err
	}
	if s.conn[udp4] != nil {
		s.conn[udp4].WriteToUDP(buf, s.cfg.IPv4Addr)
	}
	if s.conn[udp6] != nil {
		s.conn[udp6].WriteToUDP(buf, s.cfg.IPv4Addr)
	}
	return nil
}
