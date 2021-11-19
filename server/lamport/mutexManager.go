package lamport

import (
	"hotel/config"
	"math"
)

type MutexManager struct {
	SelfId   int
	H        int
	AgreedSC bool
	T        [config.NB_SERVER]MessageLamport
}

// TODO main loop and channels for receiving messages from Hotel and ConnManager
func (m MutexManager) start() {
	//for {
	//	select
	//}
}

func (m MutexManager) syncStamp(incomingStamp int) {
	m.H = int(math.Max(float64(incomingStamp), float64(m.H))) + 1 // TODO not clean ?
}

func (m MutexManager) req(msg MessageLamport) {
	m.syncStamp(msg.H)
	m.T[msg.SenderID] = msg
	if m.T[m.SelfId].Type != REQ {
		m.ack(msg)
	}
	m.verifySC()
}

func (m MutexManager) ack(msg MessageLamport) {
	m.syncStamp(msg.H)
	if m.T[msg.SenderID].Type != REQ {
		m.T[msg.SenderID] = msg
	}
	m.verifySC()
}

func (m MutexManager) rel(msg MessageLamport) {
	m.syncStamp(msg.H)
	m.T[msg.SenderID] = msg
	m.verifySC()
}

func (m MutexManager) verifySC() {
	if m.T[m.SelfId].Type != REQ {
		return
	}
	oldest := true
	for _, msg := range m.T {
		if msg.SenderID == m.SelfId {
			continue
		}
		if (m.T[m.SelfId].H > msg.H) || (m.T[m.SelfId].H == msg.H && m.SelfId > msg.SenderID) {
			oldest = false
			break
		}
	}
	if oldest {
		m.AgreedSC = true // TODO can be factorised
	}
}
