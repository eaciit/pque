# pque
Library to perform job worker pools using Golang (Go) . 

## Feature
- Configurable number of worker
- Support for user defined function of main job, post job and while waiting for completion

## Usage
```
jobCount := 1000

//--- Prepare the worker pool
q := NewQue()
q.WorkerCount = 100
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
	fmt.Printf("Processed: %d Completed: %d\n", q.ProcessedJob, q.CompletedJob)
}

//--- Sending key to indentifies job, job will run immidiately for after key rcvd
q.WaitForKeys()
for i := 0; i < jobCount; i++ {
	q.SendKey(i)
}
//--- Tell the que that all key has been sent
q.KeySendDone()

//--- Wait until all jobs done
q.WaitForCompletion()

```
