# retrying-queue
Queue implementation for retry mechanism

Usage:

```
for {
    value, err := queue.Dequeue()
    if err != nil {
        break
    }

    err := doSomething(value)
    if err != nil {
        queue.Fail()
        continue
    } 
		
    queue.Success()
}
```
