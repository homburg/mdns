package zeroconf

import (
	"log"
	"net"

	dns "github.com/miekg/godns"
)

func init() {
	listener := &listener {
		socket: openIPv4Socket(net.IPv4zero),
		zone: Local,
	}

	if err := listener.socket.JoinGroup(nil, net.IPv4(224, 0, 0, 251)); err != nil {
                log.Fatal(err)
        }

	go listener.mainloop()
}

func openIPv4Socket(ip net.IP) *net.UDPConn {
        conn, err := net.ListenUDP("udp4", &net.UDPAddr{
                IP:   ip,
                Port: 5353,
        })
        if err != nil {
                log.Fatal(err)
        }
        return conn
}

type listener struct {
	socket          *net.UDPConn
	zone	*Zone
}

func (l *listener) mainloop() {
	buf := make([]byte, 1500)
	for {
		read, _, err := l.socket.ReadFromUDP(buf[:])
		if err != nil {
			log.Fatal(err)
		}
		msg := new(dns.Msg)
		msg.Unpack(buf[:read])
		if isQuestion(msg) {
			for _, question := range msg.Question {
				for result := range l.zone.Query(question) {
					if result.publish {
					
					}
				}
			}	
		} else {
			for _, rr := range msg.Answer {
				l.zone.Add( &Entry {
					expires: 2^31,
					publish: false,
					rr: rr, 
				})
			}
		}
	}
}

func isQuestion(msg *dns.Msg) bool {
	return !msg.MsgHdr.Response
}
