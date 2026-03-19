package or

import (
	"testing"
	"time"
)

func TestOrReturnsWhenFirstChannelCloses(t *testing.T) {
	start := time.Now()
	select {
	case <-Or(sig(100*time.Millisecond),
		sig(2*time.Second),
		sig(1*time.Second)):

	case <-time.After(200 * time.Millisecond):
		t.Fatalf("Or не завершился вовремя, прошло %v", time.Since(start))
	}
}

func TestOrZeroChannels(t *testing.T) {
	res := Or()
	if res != nil {
		t.Fatalf("Or вместо nil вернул, %v", res)
	}
}

func TestOrOneChannel(t *testing.T) {
	ch := make(chan interface{})
	res := Or(ch)
	if res != ch {
		t.Fatalf("Ожидали, что Or вернет исходный канал, но получили %v", res)
	}
}

func TestOrNilAndDoneChannels(t *testing.T) {
	start := time.Now()
	var chNil <-chan interface{}
	chDone := sig(50 * time.Millisecond)
	select {
	case <-Or(chNil, chDone):
	case <-time.After(150 * time.Millisecond):
		t.Fatalf("Or не завершился вовремя, прошло %v", time.Since(start))
	}
}

func TestOrRecursiveChannels(t *testing.T) {
	start := time.Now()
	select {
	case <-Or(sig(2*time.Hour),
		sig(5*time.Minute),
		sig(10*time.Millisecond),
		sig(1*time.Hour),
		sig(10*time.Minute),
		sig(15*time.Hour),
		sig(2*time.Minute),
		sig(2*time.Second),
		sig(3*time.Hour),
		sig(1*time.Minute)):
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("Or не завершился вовремя, прошло %v", time.Since(start))
	}
}

func TestOrDoesNotCloseEarly(t *testing.T) {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})
	select {
	case <-Or(ch1, ch2, ch3):
		t.Fatalf("Or закрыл канал раньше ожидаемого")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestOrTwoNilChannels(t *testing.T) {
	var ch1Nil <-chan interface{}
	var ch2Nil <-chan interface{}
	select {
	case <-Or(ch1Nil, ch2Nil):
		t.Fatalf("Or закрыл канал")
	case <-time.After(100 * time.Millisecond):
	}
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}
