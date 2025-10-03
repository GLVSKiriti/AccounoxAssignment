# Explaination for the code in comments and answered in below questions
```
    package main

    import "fmt"

    func main() {
        // Creates a buffered channel that holds functions (func()) with capacity 10.
        cnp := make(chan func(), 10)

        // Starts 4 goroutines. Each routine runs the same code concurrently.
        for i := 0; i < 4; i++ {
            // Each goroutine continuously listens on cnp channel and executes it
            // when a function is received
            go func() {
                for f := range cnp {
                    f()
                }
            }()
        }
        // Sends a function into the channel. One of the 4 goroutines will pick it up and run it.
        cnp <- func() {
            fmt.Println("HERE1")
        }
        fmt.Println("Hello") // Prints and main function exits 
    }
```

## 1. Explaining how the highlighted constructs work?
Explained in code comments

## 2. Giving use-cases of what these constructs could be used for.
This type of code is useful in worker pools. Like example background job processing (e.g., sending emails, generating reports). \
Running many small tasks in parallel.

## 3. What is the significance of the for loop with 4 iterations?
It means it creates 4 go routines which means upto 4 tasks are executed concurrently.\
This concurrent execution improves throughput based on CPU cores we can increase number of go routines.
## 4. What is the significance of make(chan func(), 10)?
Here channel buffer size is 10 which means upto 10 tasks can be queued without blocking sender\

## 5. Why is “HERE1” not getting printed?
The program dont wait for the go routines to complete.\
main function exit too quickly so goroutine not executed and hence "HERE1" not printed.\
To remove this bug we should wait for fo routines to compelte using WaitGroup.\