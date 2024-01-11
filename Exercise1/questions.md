Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> Concurrency is the illusion of parallelism, but several tasks are not really executed at the same time. It does not finish one task before starting on another task, and that is why it seems like the tasks are executed at the same time. In reality it switches back and forth between threads.

What is the difference between a *race condition* and a *data race*? 
> A race condition happens when the order of when tasks are executed affects the results and the correctness. Data race is a type of race condition that happens when several threads access data at the same time. 
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> A scheduler has the resposibility for managing the execution of tasks. It plays an important role in achieving concurrency and multitasking by delegating the roles to the CPU.


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> by using threads we can switch between different functions, in the middle of the functions process, without loosing data. 

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> fibers are threads that you split in to several tasks, making them lighter weight and more efficient in certaint situations. 

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> this will make it more difficult to code, but the system will run more efficiently, so both

What do you think is best - *shared variables* or *message passing*?
> we think shared variables were easier to understand


