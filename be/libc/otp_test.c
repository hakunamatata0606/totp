#include "otp.h"

#include <stdio.h>
#include <time.h>
#include <string.h>
#include <unistd.h>
#include <assert.h>
#include <pthread.h>
#include <stdlib.h>

extern time_t round_time(time_t t, uint64_t interval);

const char* secret = "this is secret";
const char* client_secret = "aloha";
const int interval = 10;
const int digit = 6;

void* test(void* data)
{
    int initial_sleep = *((int*) data);
    free((int*) data);
    sleep(initial_sleep);

    char buffer1[100];
    char buffer2[100];
    time_t now = time(NULL);

    if (now + 1 >= round_time(now + interval, interval)) {
        // avoid time near interval end
        sleep(2);
    }

    uint64_t first = generate_otp(secret, sizeof(secret), client_secret, sizeof(client_secret), interval, digit);
    sleep(1);
    uint64_t second = generate_otp(secret, sizeof(secret), client_secret, sizeof(client_secret), interval, digit);
    sleep(interval);
    uint64_t third = generate_otp(secret, sizeof(secret), client_secret, sizeof(client_secret), interval, digit);

    printf("First: %lu, second: %ld, third: %lu\n", first, second, third);
    assert(first == second);
    assert(first != third);
}

int main() {
    const int num_thread = 5;
    pthread_t threads[num_thread];
    for (int i = 0; i < num_thread; i++) {
        int *initial_sleep = malloc(sizeof(int));
        *initial_sleep = i;
        pthread_create(&threads[i], NULL, test, initial_sleep);
    }

    for (int i = 0; i < num_thread; i++) {
        pthread_join(threads[i], NULL);
    }

    printf("Test passed !!!\n");
    
    return 0;
}