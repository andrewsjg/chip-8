package machine

type Input [16]bool

var waitforInput = false

func (i *Input) wait() bool {
	return waitforInput
}

func (i *Input) enableWait(enabled bool) {
	waitforInput = enabled
}
