package core

import (
	"sync"
	"sync/atomic"
)

type AppState struct {
	logSequence            int64
	currentHeight          int64
	currentTxIndex         int64
	currentCommittedHeight int64
	lock                   sync.Mutex
	restoreWg              sync.WaitGroup
}

func (s *AppState) Lock() {
	s.lock.Lock()
}

func (s *AppState) UnLock() {
	s.lock.Unlock()
}

func (s *AppState) CommitHeight() {
	atomic.AddInt64(&s.currentCommittedHeight, 1)
}

func (s *AppState) UpdateCurrentHeight() {
	atomic.AddInt64(&s.currentHeight, 1)
}

func (s *AppState) GetCurrentHeight() int64 {
	return atomic.LoadInt64(&s.currentHeight)
}

func (s *AppState) UpdateLogSequence() {
	atomic.AddInt64(&s.logSequence, 1)
}

func (s *AppState) GetLogSequence() int64 {
	return atomic.LoadInt64(&s.logSequence)
}

func (s *AppState) UpdateCurrentTxIndex() {
	atomic.AddInt64(&s.currentTxIndex, 1)
}

func (s *AppState) ClearCurrentTxIndex() {
	s.currentTxIndex = 0
}

func (s *AppState) GetCurrentTxIndex() int64 {
	return atomic.LoadInt64(&s.currentTxIndex)
}
