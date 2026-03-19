package or

import (
	"fmt"
	"time"
)

func ExampleOr() {
	ch1 := sig(10 * time.Millisecond)
	ch2 := sig(10 * time.Second)
	ch3 := sig(200 * time.Millisecond)
	<-Or(ch1, ch2, ch3)

	fmt.Println("done")
	// Output:
	// done
}
