package socket

import (
	"sync"

	"github.com/asaskevich/EventBus"
)

var bus EventBus.Bus
var wg sync.WaitGroup

func init() {
	bus = EventBus.New()
}
