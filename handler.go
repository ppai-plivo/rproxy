package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/tidwall/redcon"
)

type proxy struct {
	phase        phaseType
	srcRedisAddr string
	dstRedisAddr string
}

type proxyContext struct {
	phase phaseType
	src   net.Conn
	dst   net.Conn
}

func (p *proxy) handler(conn redcon.Conn, cmd redcon.Command) {

	ctx, ok := conn.Context().(*proxyContext)
	if !ok {
		return
	}

	command := strings.ToLower(string(cmd.Args[0]))
	if _, isCmdReadOnly := readOnlyCmds[command]; isCmdReadOnly {

		var readTarget net.Conn
		if ctx.phase == WriteBothReadSrc {
			readTarget = ctx.src
		} else {
			readTarget = ctx.dst
		}

		_, err := readTarget.Write(cmd.Raw)
		if err != nil {
			log.Printf("Write failed: %v", err)
		}

		return
	}

	if ctx.phase == WriteBothReadSrc || ctx.phase == WriteBothReadDst {
		if _, err := ctx.src.Write(cmd.Raw); err != nil {
			log.Printf("src.Write failed: %v", err)
		}
	}

	if _, err := ctx.dst.Write(cmd.Raw); err != nil {
		log.Printf("dst.Write failed: %v", err)
	}
}

func (p *proxy) relayReplies(client io.Writer, server io.Reader) {
	for {
		_, err := io.Copy(client, server)
		if err == io.EOF {
			return
		}
		if checkNetOpError(err) != nil {
			log.Printf("io.Copy error: %v", err)
			return
		}
	}
}

// onAccept is called when a client connects. If the function
// returns true, connection is accepted.
func (p *proxy) onAccept(conn redcon.Conn) bool {
	log.Printf("client connected: %s\n", conn.RemoteAddr())

	src, err := net.Dial("tcp", p.srcRedisAddr)
	if err != nil {
		log.Printf("net.Dial(%s) failed: %v", p.srcRedisAddr, err)
		return false
	}

	dst, err := net.Dial("tcp", p.dstRedisAddr)
	if err != nil {
		log.Printf("net.Dial(%s) failed: %v", p.dstRedisAddr, err)
		src.Close()
		return false
	}

	switch p.phase {
	case WriteBothReadSrc:
		go p.relayReplies(conn.NetConn(), src)
		go p.relayReplies(ioutil.Discard, dst)
	case WriteBothReadDst:
		go p.relayReplies(conn.NetConn(), dst)
		go p.relayReplies(ioutil.Discard, src)
	case WriteDstReadDst:
		go p.relayReplies(conn.NetConn(), dst)
	}

	conn.SetContext(&proxyContext{
		src:   src,
		dst:   dst,
		phase: p.phase,
	})

	return true
}

// onClose is called when a client connection is disconnected.
func (p *proxy) onClose(conn redcon.Conn, err error) {
	log.Printf("client disconnected: %s, err: %v\n", conn.RemoteAddr(), err)

	ctx, ok := conn.Context().(*proxyContext)
	if ok {
		ctx.src.Close()
		ctx.dst.Close()
	}
}

func checkNetOpError(err error) error {
	if err != nil {
		netOpError, ok := err.(*net.OpError)
		if ok && strings.HasSuffix(netOpError.Err.Error(), "use of closed network connection") {
			return nil
		}
	}
	return err
}
