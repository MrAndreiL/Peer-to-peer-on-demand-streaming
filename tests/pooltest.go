package tests

import (
	"Peer-to-peer-on-demand-streaming/utils"
	"testing"
)

func TestThreadPool_CreatePool(t *testing.T) {
	if _, err := utils.NewThreadPool(0, 0); err == utils.ErrNoThreads {
		t.Fatalf("expected error when creating pool with 0 workers, got: %v", err)
	}
	if _, err := utils.NewThreadPool(-1, 0); err == utils.ErrNoThreads {
		t.Fatalf("expected error when creating pool with -1 workers, got: %v", err)
	}
	if _, err := utils.NewThreadPool(1, -1); err == utils.ErrNegativeChannelSize {
		t.Fatalf("expected error when creating pool with -1 channel size, got: %v", err)
	}
}

func TestThreadPool_PanicStartStop(t *testing.T) {
	pool, err := utils.NewThreadPool(5, 1)
	if err != nil {
		t.Fatal("Error creating pool:", err)
	}

	pool.Start()
	pool.Start()

	pool.Stop()
	pool.Stop()
}
