package main

import (
	"fmt"
	"workers/threads"
)

func jobProcessor(input interface{}) error {
	fmt.Println("Job processing", input)
	return nil
}

func outputProcessor(ouput threads.Output) error {
	fmt.Println("Processing output for job", ouput.Job.ID)
	return nil
}

func main() {
	workerPool := threads.NewWorkerPool(3)
	inputStrings := []string{"test1", "test2", "test3"}
	inputs := make([]interface{}, len(inputStrings))
	for index, str := range inputStrings {
		inputs[index] = str
	}
	workerPool.RunJobs(inputs, jobProcessor, outputProcessor)
}
