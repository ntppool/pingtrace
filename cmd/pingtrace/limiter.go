package main

type Limiter struct {
	runSemaphore chan int
}

func NewLimiter(length int) *Limiter {
	return &Limiter{
		runSemaphore: make(chan int, 10),
	}
}

func (l *Limiter) Check() bool {
	select {
	case l.runSemaphore <- 1:
		<-l.runSemaphore
		return true
	default:
		return false
	}
}

func (l *Limiter) Allowed() bool {
	select {
	case l.runSemaphore <- 1:
		return true
	default:
		return false
	}
}

func (l *Limiter) Done() {
	<-l.runSemaphore
}
