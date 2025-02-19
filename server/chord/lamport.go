package chord

type LamportClock struct {
	time int
}

func (l *LamportClock) Tick() {
	l.time++
}

func (l *LamportClock) Update(time int) {
	l.time = time
}

func (l *LamportClock) Time() int {
	return l.time
}
