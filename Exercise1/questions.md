Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
> Concurrency is the illusion of parallelism, but several tasks are not really executed at the same time. It does not finish one task before starting on another task, and that is why it seems like the tasks are executed at the same time. In reality it switches back and forth between threads.

What is the difference between a *race condition* and a *data race*? 
> A race condition happens when the order of when tasks are executed affects the results and the correctness. Data race 
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> *Your answer here* 


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> *Your answer here*

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> *Your answer here*

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> *Your answer here*

What do you think is best - *shared variables* or *message passing*?
> *Your answer here*


