package threads

import (
	"fmt"
	"math/rand"
	"sync"
)

// Job for processing in pool
type Job struct {
	ID    int
	Input interface{}
}

// Output is the result of processing a job
type Output struct {
	Job Job
	Err error
}

// WorkerPool is a data structure to hold a pool of workers
type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	outputs    chan Output
	done       chan bool
}

// JobProcessor is a function used for processing individual jobs
type JobProcessor func(interface{}) error

// OutputProcessor is a function for processing outputs
type OutputProcessor func(Output) error

// NewWorkerPool returns a worker pool allocated with numWorkers
func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := WorkerPool{numWorkers: numWorkers}
	pool.jobs = make(chan Job, numWorkers)
	pool.outputs = make(chan Output, numWorkers)
	return &pool
}

// RunJobs runs the submitted inputs
func (w *WorkerPool) RunJobs(inputs []interface{}, jobProcessor JobProcessor, outputProcessor OutputProcessor) {
	w.done = make(chan bool)
	go w.submitJobs(inputs)
	go w.collectOutputs(outputProcessor)
	go w.processJobs(jobProcessor)
	<-w.done
}

func (w *WorkerPool) submitJobs(inputs []interface{}) {
	for input := range inputs {
		id := rand.Intn(len(inputs) * 1000)
		job := Job{id, input}
		w.jobs <- job
	}
	close(w.jobs)
}

func (w *WorkerPool) collectOutputs(outputProcessor OutputProcessor) {
	for output := range w.outputs {
		err := outputProcessor(output)
		if err != nil {
			fmt.Println("Failed to process output for job ", output.Job.ID)
		}
	}
	w.done <- true
}

func (w *WorkerPool) processJobs(jobProcessor JobProcessor) {
	var wg sync.WaitGroup
	for i := 0; i < w.numWorkers; i++ {
		wg.Add(1)
		go w.process(&wg, jobProcessor)
	}
	wg.Wait()
	close(w.outputs)
}

func (w *WorkerPool) process(wg *sync.WaitGroup, jobProcessor JobProcessor) {
	for job := range w.jobs {
		result := jobProcessor(job.Input)
		output := Output{job, result}
		w.outputs <- output
		if result != nil {
			fmt.Println("Failed to process job for Id", job.ID)
		}
	}
	wg.Done()
}
