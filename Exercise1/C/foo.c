// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t mtx;

// Note the return type: void*
void* incrementingThreadFunction(){
    // TODO: increment i 1_000_000 times
    int a;
    for (a = 0; a < 1000000; a++){
        pthread_mutex_lock (&mtx);
        i++;
        pthread_mutex_unlock (&mtx);
    }

    return NULL;
}

void* decrementingThreadFunction(){
    // TODO: decrement i 1_000_000 times
    int b;
    for (b = 0; b < 1000000; b++){
        pthread_mutex_lock (&mtx);
        i--;
        pthread_mutex_unlock (&mtx);
    }
    return NULL;
}


int main(){
    //MUTEX

    //2nd arg is a pthread_mutexattar_t
    pthread_mutex_init (&mtx, NULL);

    // TODO: 
    // Declare thread identifiers
    pthread_t Thread1;
    pthread_t Thread2;
    
    // TODO:
    // Create threads using `pthread_create`
    pthread_create(&Thread1, NULL, incrementingThreadFunction, NULL);
    pthread_create(&Thread2, NULL, decrementingThreadFunction, NULL);
    
    // TODO:
    // Wait for the two threads to be done before printing the final result
    pthread_join(Thread1, NULL);
    pthread_join(Thread2, NULL);

    pthread_mutex_destroy (&mtx);


    printf("The magic number is: %d\n", i);
    return 0;
}
