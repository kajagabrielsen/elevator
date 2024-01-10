The result, i, is not equal to zero and the answere changes every time.
The reason why is because they do not run simoltaniosly.
To resolve the problem we need to use one of the two;
    - Intruduce a botleneck
    - Send messanges instead of sharing memory
We use Mutex instead of Semaphore becuase with Semaphore anyone can have the key, but Mutex rquires that ony the one that lockes has the key. 
