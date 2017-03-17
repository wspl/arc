package arc

type LoopTaskDelegate func() (bool, error)

func NewLoopTask(delegate LoopTaskDelegate) *LoopTask {
	return &LoopTask{
		Delegate: delegate,
	}
}

type LoopTask struct {
	Enabled  bool
	Delegate LoopTaskDelegate
	OnError func(error)
}

func (t *LoopTask) Start() *LoopTask {
	t.Enabled = true
	go func() {
		for t.Enabled {
			keep, err := t.Delegate()
			if err != nil && t.OnError != nil {
				t.OnError(err)
			}
			if !keep {
				t.Enabled = false
				break
			}
		}
	}()
	return t
}

func (t *LoopTask) Stop() {
	t.Enabled = false
}