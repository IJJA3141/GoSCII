package tui

type CheckBox struct {
	checked bool
}

func (bx *CheckBox) HandleKey(key string) (consumed bool, submit bool) {
	if key == KEY_ENTER {
		bx.checked = !bx.checked
		return true, true
	}

	return false, false
}
