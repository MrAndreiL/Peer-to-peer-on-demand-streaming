package main

import "Peer-to-peer-on-demand-streaming/utils"

func OpenThreadPool() utils.Pool {
	pool := utils.CreateThreadPool(utils.Threads, utils.Tasks)
	pool.Start()
	return pool
}
