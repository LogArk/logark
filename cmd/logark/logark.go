package main

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/LogArk/logark/internal/filters/mutate"
	"github.com/LogArk/logark/internal/outputs/stdout"
	"github.com/LogArk/logark/internal/pipeline"
	"github.com/LogArk/logark/internal/queue"
)

func execPipeline(log *gabs.Container, p pipeline.Pipeline) {
	var executionStack []pipeline.Filter

	// Initialize filter stack
	for i := len(p.Process) - 1; i >= 0; i-- {
		executionStack = append(executionStack, p.Process[i])
	}

	for len(executionStack) > 0 {
		n := len(executionStack) - 1

		// Pop last element
		V := executionStack[n]
		executionStack = executionStack[:n]

		//fmt.Println("----", V.GetName(), "----")

		status := false
		switch V.GetName() {
		case "mutate":
			status = mutate.ExecFilter(log, V.GetParams())
		default:
			//fmt.Println("Cannot handle", V.GetName())
		}
		if status {
			//fmt.Println("Operation was succcesful, looking at on_success")
			success := V.GetOnSuccess()
			for i := len(success) - 1; i >= 0; i-- {
				executionStack = append(executionStack, success[i])
			}
		} else {
			//fmt.Println("Operation was not succcesful, looking at on_failure")
			failure := V.GetOnFailure()
			for i := len(failure) - 1; i >= 0; i-- {
				executionStack = append(executionStack, failure[i])
			}
		}

	}
}

func execOutput(log []byte, p pipeline.Pipeline) {
	for _, o := range p.Outputs {
		switch o.GetName() {
		case "stdout":
			stdout.Send(log)
		default:
		}
	}
}

func filterWorker(qm *queue.QueueManager, p pipeline.Pipeline, workerId uint) {
	for {
		job, _ := qm.GetFilterJob()
		//fmt.Println(workerId, " : got filter job:  ", job.JobId)
		jsonParsed, _ := gabs.ParseJSON(job.Log)
		execPipeline(jsonParsed, p)
		job.Log = jsonParsed.Bytes()
		qm.CompleteFilterJob(job)
	}
}

func outputWorker(qm *queue.QueueManager, p pipeline.Pipeline, workerId uint) {
	for {
		job, _ := qm.GetOutputJob()
		//fmt.Println(workerId, " : got output job:  ", job.JobId)
		execOutput(job.Log, p)
		qm.CompleteOutputJob(job)
	}
}

func main() {

	fmt.Println("Loading pipeline...")
	p, _ := pipeline.Load("./config/pipeline.yaml")

	fmt.Println("Creating queue manager")
	qm := queue.NewQueueManager()

	go outputWorker(qm, p, 0)

	for i := uint(0); i < p.Settings.Workers; i++ {
		fmt.Println("Starting worker: ", i)
		go filterWorker(qm, p, i)
	}

	fmt.Println("Starting Output dispatch")
	go qm.OutputDispatch()

	fmt.Println("Starting Filter dispatch")
	go qm.FilterDispatch()

	for {
		var input string
		fmt.Scanln(&input)
		qm.PushLog([]byte(input))
	}

	// Input
	return
	/*
		jsonParsed, _ := gabs.ParseJSON(e)

		execPipeline(jsonParsed, p)

		// Output
		fmt.Println(jsonParsed.String())
	*/
}
