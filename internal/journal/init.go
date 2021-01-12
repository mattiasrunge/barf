package journal

import (
	"os"
	"path"

	"rft/internal/com/server"
	"rft/internal/config"

	"github.com/asaskevich/EventBus"
)

var bus EventBus.Bus
var activeDir = ""
var historyDir = ""

func init() {
	bus = EventBus.New()

	activeDir = path.Join(config.JournalDir, "active")
	historyDir = path.Join(config.JournalDir, "history")

	os.Mkdir(activeDir, 0700)
	os.Mkdir(historyDir, 0700)

	server.OnOperationCreate(create)
	server.OnOperationAbort(abort)
	server.OnOperationStatus(status)
	server.OnListOperations(list)
}
