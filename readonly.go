package main

import (
	"github.com/go-redis/redis"
)

// a set containing list of read-only ommands
var readOnlyCmds = make(map[string]struct{})

func getReadOnlyCommands(addr string) error {

	redisdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{addr},
	})
	defer redisdb.Close()

	if _, err := redisdb.Ping().Result(); err != nil {
		return err
	}

	cmds, err := redisdb.Command().Result()
	if err != nil {
		return err
	}

	for cmd, info := range cmds {
		if info.ReadOnly {
			readOnlyCmds[cmd] = struct{}{}
		}
	}

	return nil
}
