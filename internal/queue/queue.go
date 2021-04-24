package queue

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type QueueManager struct {
	preFilterQueue      []Job
	preFilterQueueMutex sync.Mutex

	postFilterQueue      []Job
	postFilterQueueMutex sync.Mutex

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

		qm.preFilterQueueMutex.Lock()
		for i, v := range qm.preFilterQueue {
			if v.isLocked == false {
				qm.preFilterQueue[i].isLocked = true
				j = qm.preFilterQueue[i]
				found = true
				break
			}
		}
		qm.preFilterQueueMutex.Unlock()

		if found {
			fmt.Println("Dispatching job", j.JobId)
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

		qm.postFilterQueueMutex.Lock()
		for i, v := range qm.postFilterQueue {
			if !v.isLocked {
				qm.postFilterQueue[i].isLocked = true
				j = qm.postFilterQueue[i]
				found = true
				break
			}
		}
		qm.postFilterQueueMutex.Unlock()

		if found {
			fmt.Println("Dispatching output", j.JobId)
			qm.outputDispatch <- j
		} else {
			time.Sleep(time.Second)
		}

	}
}

func (qm *QueueManager) PushLog(log []byte) error {
	var j Job
	qm.preFilterQueueMutex.Lock()

	j.JobId = qm.jobCounter
	j.Log = log
	j.isLocked = false
	qm.preFilterQueue = append(qm.preFilterQueue, j)
	qm.jobCounter++

	qm.preFilterQueueMutex.Unlock()
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
	return len(qm.preFilterQueue)
}

func (qm *QueueManager) CompleteFilterJob(job Job) error {
	prefilterIndex := 0
	found := false

	// Find job in prefilter
	for i, v := range qm.preFilterQueue {
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

	// Queue job in postfilter
	qm.postFilterQueueMutex.Lock()
	job.isLocked = false
	qm.postFilterQueue = append(qm.postFilterQueue, job)
	qm.postFilterQueueMutex.Unlock()

	// Remove job from prefilter
	qm.preFilterQueueMutex.Lock()
	qm.preFilterQueue = append(qm.preFilterQueue[:prefilterIndex], qm.preFilterQueue[prefilterIndex+1:]...)
	qm.preFilterQueueMutex.Unlock()

	return nil
}

func (qm QueueManager) Dump() {
	fmt.Println("preFilter")
	for _, v := range qm.preFilterQueue {
		fmt.Println(v)
	}

	fmt.Println("postFilter")
	for _, v := range qm.postFilterQueue {
		fmt.Println(v)
	}
}
