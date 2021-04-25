package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Jeffail/gabs"
	"github.com/LogArk/logark/internal/outputs/stdout"
	"github.com/LogArk/logark/internal/pipeline"
	"github.com/LogArk/logark/internal/queue"
)

func execPipeline(log *gabs.Container, p *pipeline.FilterPipeline) {
	var executionStack []pipeline.FilterAction

	// Initialize filter stack
	for i := len(p.Filters) - 1; i >= 0; i-- {
		executionStack = append(executionStack, p.Filters[i])
	}

	for len(executionStack) > 0 {
		n := len(executionStack) - 1

		// Pop last element
		V := executionStack[n]
		executionStack = executionStack[:n]

		status := false
		status, _ = V.Plugin.Exec(log)

		if status {
			success := V.OnSuccess
			for i := len(success.Filters) - 1; i >= 0; i-- {
				executionStack = append(executionStack, success.Filters[i])
			}
		} else {
			//fmt.Println("Operation was not succcesful, looking at on_failure")
			failure := V.OnFailure
			for i := len(failure.Filters) - 1; i >= 0; i-- {
				executionStack = append(executionStack, failure.Filters[i])
			}
		}
	}
}

func execOutput(log []byte, o []pipeline.Output) {
	for _, o := range o {
		switch o.GetName() {
		case "stdout":
			stdout.Send(log)
		default:
		}
	}
}

func filterWorker(qm *queue.QueueManager, p *pipeline.FilterPipeline, workerId uint) {
	for {
		job, _ := qm.GetFilterJob()
		jsonParsed, err := gabs.ParseJSON(job.Log)
		if err != nil {
			fmt.Println("====================", string(job.Log))
			fmt.Println(err)
		} else {
			execPipeline(jsonParsed, p)
			job.Log = jsonParsed.Bytes()
		}

		qm.CompleteFilterJob(job)
	}
}

func outputWorker(qm *queue.QueueManager, o []pipeline.Output, workerId uint) {
	for {
		job, _ := qm.GetOutputJob()
		execOutput(job.Log, o)
		qm.CompleteOutputJob(job)
	}
}

func main() {

	p, _ := pipeline.Load("./config/pipeline.yaml")

	qm := queue.NewQueueManager()

	go outputWorker(qm, p.Outputs, 0)
	for i := uint(0); i < p.Settings.Workers; i++ {
		go filterWorker(qm, p.FilterPipeline, i)
	}
	go qm.OutputDispatch()
	go qm.FilterDispatch()

	/* Fake stdin input */
	scanner := bufio.NewScanner(os.Stdin)
	for {
		for scanner.Scan() {
			text := scanner.Text()
			qm.PushLog([]byte(text))
		}
	}
}
