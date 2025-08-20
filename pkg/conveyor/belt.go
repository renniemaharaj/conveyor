package conveyor

type ConveyorBelt struct {
	C chan Job
}

func NewConveyorBelt() *ConveyorBelt {
	return &ConveyorBelt{C: make(chan Job, 100)}
}

func (b *ConveyorBelt) Push(j *Job) {
	b.C <- *j
}

func (b *ConveyorBelt) Take() *Job {
	j := <-b.C
	return &j
}
