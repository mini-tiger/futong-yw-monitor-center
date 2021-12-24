package hub

type Hub struct {
}

var MyHub = &Hub{}

func (h *Hub) Run() {
	//go sendHostInfo()

	//go sendHostMetrics()
	sendHostMetrics()

	//go findScriptTask()

	//for {
	//	select {
	//	case ret := <-h.ScriptTask.Ch:
	//		sendScriptResult(ret)
	//	default:
	//		time.Sleep(5*time.Second)
	//	}
	//}
}
