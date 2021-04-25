package queue

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type QueueManager struct {
	filterQueue      []Job
	filterQueueMutex sync.Mutex

	outputQueue      []Job
	outputQueueMutex sync.Mutex

	jobCounter uint64

	filterDispatch chan Job
	outputDispatch chan Job
}

type Job struct {
	JobId    uint64
	Log      []byte
	isLocked bool
}

func NewQueueManager() *QueueManager {
	var newQM QueueManager
	newQM.filterDispatch = make(chan Job)
	newQM.outputDispatch = make(chan Job)
	return &newQM
}

func (qm *QueueManager) FilterDispatch() {
	for true {
		// Find next free job
		var j Job
		found := false

		qm.filterQueueMutex.Lock()
		for i, v := range qm.filterQueue {
			if v.isLocked == false {
				qm.filterQueue[i].isLocked = true
				j = qm.filterQueue[i]
				found = true
				break
			}
		}
		qm.filterQueueMutex.Unlock()

		if found {
			//fmt.Println("Dispatching job", j.JobId)
			qm.filterDispatch <- j
		} else {
			time.Sleep(time.Second)
		}

	}
}

func (qm *QueueManager) OutputDispatch() {
	for {
		// Find next free job
		var j Job
		found := false

		qm.outputQueueMutex.Lock()
		for i, v := range qm.outputQueue {
			if !v.isLocked {
				qm.outputQueue[i].isLocked = true
				j = qm.outputQueue[i]
				found = true
				break
			}
		}
		qm.outputQueueMutex.Unlock()

		if found {
			//fmt.Println("Dispatching output", j.JobId)
			qm.outputDispatch <- j
		} else {
			time.Sleep(time.Second)
		}

	}
}

func (qm *QueueManager) PushLog(log []byte) error {
	var j Job
	qm.filterQueueMutex.Lock()

	j.JobId = qm.jobCounter
	j.Log = log
	j.isLocked = false
	qm.filterQueue = append(qm.filterQueue, j)
	qm.jobCounter++

	qm.filterQueueMutex.Unlock()
	return nil
}

func (qm *QueueManager) GetFilterJob() (Job, error) {
	j := <-qm.filterDispatch
	return j, nil
}

func (qm *QueueManager) GetOutputJob() (Job, error) {
	j := <-qm.outputDispatch
	return j, nil
}

func (qm *QueueManager) GetDepth() int {
	return len(qm.filterQueue)
}

func (qm *QueueManager) CompleteFilterJob(job Job) error {
	prefilterIndex := 0
	found := false

	// Queue job in postfilter
	qm.outputQueueMutex.Lock()
	job.isLocked = false
	qm.outputQueue = append(qm.outputQueue, job)
	qm.outputQueueMutex.Unlock()

	// Remove job from prefilter
	qm.filterQueueMutex.Lock()
	// Find job in prefilter
	for i, v := range qm.filterQueue {
		if v.JobId == job.JobId {
			prefilterIndex = i
			found = true
			break
		}
	}
	if !found {
		fmt.Errorf("Job not found in prefilter queue, this is really unexpected \n")
		return errors.New("Missing job in prefilter queue")
	}
	qm.filterQueue = append(qm.filterQueue[:prefilterIndex], qm.filterQueue[prefilterIndex+1:]...)
	qm.filterQueueMutex.Unlock()

	return nil
}

func (qm *QueueManager) CompleteOutputJob(job Job) error {
	outputIndex := 0
	found := false

	// Remove job from output
	qm.outputQueueMutex.Lock()
	// Find job in outpout
	for i, v := range qm.outputQueue {
		if v.JobId == job.JobId {
			outputIndex = i
			found = true
			break
		}
	}
	if !found {
		fmt.Errorf("Job not found in output queue, this is really unexpected \n")
		return errors.New("Missing job in output queue")
	}
	qm.outputQueue = append(qm.outputQueue[:outputIndex], qm.outputQueue[outputIndex+1:]...)
	qm.outputQueueMutex.Unlock()

	return nil
}

func (qm QueueManager) Dump() {
	fmt.Println("preFilter")
	for _, v := range qm.filterQueue {
		fmt.Println(v)
	}

	fmt.Println("postFilter")
	for _, v := range qm.outputQueue {
		fmt.Println(v)
	}
}
