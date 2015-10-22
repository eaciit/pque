package pque

import (
	"fmt"
	"testing"
	"time"
)

func Test_Que(t *testing.T) {
	fmt.Println("Test Que")
	jobCount := 5

	result := make([]int, 0)

	q := NewQue()
	q.WorkerCount = 2
	q.Fn = func(in interface{}) interface{} {
		i := in.(int)
		fmt.Printf("Rcvd key: %d\n", i)
		time.Sleep(2 * time.Millisecond)
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

	if len(result) == jobCount {
		fmt.Printf("Number of result is %d, value are %v \n", len(result), result)
	} else {
		t.Errorf("Out of %d jobs only %d has been completed, value are %v", jobCount, len(result), result)
	}
}
