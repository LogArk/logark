package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/Jeffail/gabs"
	"github.com/LogArk/logark/internal/outputs/stdout"
	"github.com/LogArk/logark/internal/pipeline"
	"github.com/LogArk/logark/internal/queue"
	plg "github.com/LogArk/logark/pkg/plugin"
)

type FilterAction struct {
	Name      string
	Plugin    plg.FilterPlugin
	OnSuccess *FilterPipeline
	OnFailure *FilterPipeline
}

type FilterPipeline struct {
	Filters []FilterAction
}

func buildPipeline(f []pipeline.Filter) *FilterPipeline {
	var err error
	var rootPipeline FilterPipeline
	var plugins map[string]*plugin.Plugin = make(map[string]*plugin.Plugin)

	/* Get binary directory */
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, v := range f {
		var fa FilterAction
		fa.Name = v.GetName()

		// Check if plugin is already loaded. If not, load it
		plug := plugins[fa.Name]
		if plug == nil {
			plug, err = plugin.Open(dir + "/plugins/" + fa.Name + ".so")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			plugins[fa.Name] = plug
		}

		newPlugin, err := plug.Lookup("New")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fa.Plugin = newPlugin.(func() plg.FilterPlugin)()
		fa.Plugin.Init(v.GetParams())
		fa.OnSuccess = buildPipeline(v.GetOnSuccess())
		fa.OnFailure = buildPipeline(v.GetOnFailure())

		rootPipeline.Filters = append(rootPipeline.Filters, fa)
	}
	return &rootPipeline
}

func execPipeline(log *gabs.Container, p *FilterPipeline) {
	var executionStack []FilterAction

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

func execOutput(log []byte, p pipeline.Pipeline) {
	for _, o := range p.Outputs {
		switch o.GetName() {
		case "stdout":
			stdout.Send(log)
		default:
		}
	}
}

func filterWorker(qm *queue.QueueManager, p *FilterPipeline, workerId uint) {
	for {
		job, _ := qm.GetFilterJob()
		//fmt.Println(workerId, " : got filter job:  ", job.JobId)
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

func outputWorker(qm *queue.QueueManager, p pipeline.Pipeline, workerId uint) {
	for {
		job, _ := qm.GetOutputJob()
		//fmt.Println(workerId, " : got output job:  ", job.JobId)
		execOutput(job.Log, p)
		qm.CompleteOutputJob(job)
	}
}

func main() {

	//fmt.Println("Loading pipeline...")
	p, _ := pipeline.Load("./config/pipeline.yaml")

	rp := buildPipeline(p.Process)

	fmt.Println(rp)

	//fmt.Println("Creating queue manager")
	qm := queue.NewQueueManager()

	go outputWorker(qm, p, 0)

	for i := uint(0); i < p.Settings.Workers; i++ {
		//fmt.Println("Starting worker: ", i)
		go filterWorker(qm, rp, i)
	}

	//fmt.Println("Starting Output dispatch")
	go qm.OutputDispatch()

	//fmt.Println("Starting Filter dispatch")
	go qm.FilterDispatch()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		for scanner.Scan() {
			text := scanner.Text()
			qm.PushLog([]byte(text))
		}
	}
}
