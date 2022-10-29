package input

type MultiHandler struct {
	list []*Handler
}

func (h *MultiHandler) AddHandler(handler *Handler) {
	h.list = append(h.list, handler)
}

func (h *MultiHandler) ActionIsJustPressed(action Action) bool {
	for i := range h.list {
		if h.list[i].ActionIsJustPressed(action) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) ActionIsPressed(action Action) bool {
	for i := range h.list {
		if h.list[i].ActionIsPressed(action) {
			return true
		}
	}
	return false
}
