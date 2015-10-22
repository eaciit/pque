package pque

import (
	"fmt"
	"testing"
	"time"
)

func Test_Que(t *testing.T) {
	fmt.Println("Test Que")
	jobCount := 10

	result := make([]int, 0)

	q := NewQue()
	q.WorkerCount = 20
	q.Fn = func(in interface{}) interface{} {
		i := in.(int)
		fmt.Printf("Rcvd key: %d\n", i)
		time.Sleep(2 * time.Second)
		return i
	}
	q.FnDone = func(in interface{}) {
		i := in.(int)
		fmt.Printf("Finish processing: %d\n", i)
		result = append(result, i)
	}
	q.FnWaiting = func(q *Que) {
		for q.Completed == false {
			select {
			case <-time.After(1 * time.Second):
				fmt.Printf("Processed: %d Completed: %d\n", q.ProcessedJob, q.CompletedJob)
			}
		}
	}

	q.WaitForKeys()
	for i := 0; i < jobCount; i++ {
		q.SendKey(i)
	}
	q.KeySendDone()

	q.WaitForCompletion()
	fmt.Printf("Number of result is %d, value are %v \n", len(result), result)
}
