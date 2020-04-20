package mdns

import (
	"fmt"
	"github.com/miekg/dns"
	"go.uber.org/atomic"
	"log"
	"net"
	"strings"
)

const forceUnicastResponses = false

// Server ...
type Server interface {
}

type server struct {
	serviceAddr  string
	instanceAddr string
	enumAddr     string
	cfg          *OptionConfig
	conn         []*net.UDPConn
	stop         *atomic.Bool
}

// Start ...
func (s *server) Start() {
	for i := range s.conn {
		if s.conn[i] != nil {

		}
	}
}

// recv is a long running routine to receive packets from an interface
func (s *server) recv(c *net.UDPConn) {
	buf := make([]byte, 65536)
	for !s.stop.Load() {
		n, from, err := c.ReadFrom(buf)
		if err != nil {
			continue
		}
		if err := s.parsePacket(buf[:n], from); err != nil {
			logE("failed to handle query", "error", err)
		}
	}
}

// parsePacket is used to parse an incoming packet
func (s *server) parsePacket(packet []byte, from net.Addr) error {
	var msg dns.Msg
	if err := msg.Unpack(packet); err != nil {
		log.Printf("[ERR] mdns: Failed to unpack packet: %v", err)
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

	if s.cfg.LogEmptyResponses && len(multicastAnswer) == 0 && len(unicastAnswer) == 0 {
		questions := make([]string, len(query.Question))
		for i, q := range query.Question {
			questions[i] = q.Name
		}
		log.Printf("no responses for query with questions: %s", strings.Join(questions, ", "))
	}

	if mresp := resp(false); mresp != nil {
		if err := s.sendResponse(mresp, from, false); err != nil {
			return fmt.Errorf("mdns: error sending multicast response: %v", err)
		}
	}
	if uresp := resp(true); uresp != nil {
		if err := s.sendResponse(uresp, from, true); err != nil {
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
	records := s.zoneInstance(q)

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

func (s *server) zoneInstance(q dns.Question) []dns.RR {
	switch q.Name {
	case s.cfg.EnumAddr:
		return s.enumRecords(q)
	case s.cfg.ServiceAddr:
		return s.serviceRecords(q)
	case s.cfg.InstanceAddr:
		return s.instanceRecords(q)
	case s.cfg.HostName:
		if q.Qtype == dns.TypeA || q.Qtype == dns.TypeAAAA {
			return s.instanceRecords(q)
		}
		fallthrough
	default:
		return nil
	}
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
			Ptr: s.cfg.ServiceAddr,
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
			Ptr: s.cfg.InstanceAddr,
		}
		servRec := []dns.RR{rr}

		// Get the instance records
		instRecs := s.instanceRecords(dns.Question{
			Name:  s.cfg.InstanceAddr,
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
			Name:  s.cfg.InstanceAddr,
			Qtype: dns.TypeSRV,
		})

		// Add the TXT record
		recs = append(recs, s.instanceRecords(dns.Question{
			Name:  s.cfg.InstanceAddr,
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
			Name:  s.cfg.InstanceAddr,
			Qtype: dns.TypeA,
		})...)

		// Add the AAAA record
		recs = append(recs, s.instanceRecords(dns.Question{
			Name:  s.cfg.InstanceAddr,
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