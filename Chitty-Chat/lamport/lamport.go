package lamport

type LamportTime struct {
	time   uint32
	Client string
}

func (time *LamportTime) Increment() {
	time.time = time.time + 1
	//log.Printf("Current process timestamp:  %d", time.time)
}

func (time *LamportTime) GetTimestamp() uint32 {
	return time.time
}
