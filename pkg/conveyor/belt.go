package conveyor

// The channel, conveyor belt for all jobs
var CONVEYOR_BELT = make(chan Job, 100)
