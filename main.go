package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tidwall/redcon"
)

func main() {

	listenAddr := flag.String("addr", ":6379",
		"Address on which rproxy will listen on")
	srcRedisAddr := flag.String("src", "",
		"Address of source redis: Example: redis-nonclusterdev.example.com:6379")
	dstRedisAddr := flag.String("dst", "",
		"Address of destination redis: Example: redis-clusterdev.example.com:6379")
	migrationPhase := flag.Int("phase", int(WriteBothReadSrc), fmt.Sprintf(`Migration phase. Possible values:
	(%d: %s)(default)
	(%d: %s)
	(%d: %s)`,
		WriteBothReadSrc, WriteBothReadSrc.String(),
		WriteBothReadDst, WriteBothReadDst.String(),
		WriteDstReadDst, WriteDstReadDst.String()))
	flag.Parse()

	if *srcRedisAddr == "" || *dstRedisAddr == "" {
		log.Fatal("src and dst addrs cannot be empty")
	}

	if err := getReadOnlyCommands(*srcRedisAddr); err != nil {
		log.Fatal(err)
	}
	log.Printf("Source redis reachable at %s\n", *srcRedisAddr)

	if err := getReadOnlyCommands(*dstRedisAddr); err != nil {
		log.Fatal(err)
	}
	log.Printf("Destination redis reachable at %s\n", *dstRedisAddr)

	phase := phaseType(*migrationPhase)
	if phase != WriteBothReadSrc && phase != WriteBothReadDst && phase != WriteDstReadDst {
		log.Fatalf("Invalid migration phase %d", *migrationPhase)
	}
	log.Printf("Chosen migration phase %d: %s\n", phase, phase.String())

	p := &proxy{
		phase:        phase,
		srcRedisAddr: *srcRedisAddr,
		dstRedisAddr: *dstRedisAddr,
	}

	log.Printf("Serving at %s\n", *listenAddr)
	if err := redcon.ListenAndServe(*listenAddr, p.handler, p.onAccept, p.onClose); err != nil {
		log.Fatal(err)
	}
}
