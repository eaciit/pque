package pque

import (
	"fmt"
	"testing"
)

func Test_Que(t *testing.T) {
	fmt.Println("Test Que")
	jobCount := 2

	q := NewQue()
	q.WorkerCount = 2
	q.WaitForKeys()

	for i := 0; i < jobCount; i++ {
		q.SendKey(i)
	}
	q.KeySendDone()

	q.WaitForCompletion()
}
