#include "otp.h"

#include <time.h>
#include <string.h>
#include <unistd.h>
#include <assert.h>
#include <math.h>
#include <openssl/hmac.h>
#include <openssl/err.h>

const char* time_format = "%Y-%m-%d %H:%M:%S";

static size_t get_time_string(time_t t, char* buffer, size_t buffer_size)
{
    struct tm *timeinfo;
    memset(buffer, 0, buffer_size);

    timeinfo = localtime(&t);
    return strftime(buffer, buffer_size, time_format, timeinfo);
}

time_t round_time(time_t t, uint64_t interval)
{
    return ((t + interval / 2) / interval) * interval;
}

static int generate_hash(char* output, const char* secret, size_t secret_size, const char* input, size_t input_size)
{
    // using openssl > 1.1.0 no need init library
    unsigned int len;
    unsigned char* ret = HMAC(EVP_sha256(), secret, secret_size, input, input_size, (unsigned char*) output, &len);
    if (!ret) {
        int error_code = ERR_get_error();
        char* error = ERR_error_string(error_code, NULL);
        printf("generate_hash(): Failed to generate hash: %s\n", error);
        return error_code;
    }
    return 0;
}

uint64_t generate_otp(const char* secret, size_t secret_size, const char* client_secret, size_t client_secret_size, uint64_t interval, uint8_t digit)
{
    time_t time_rounded;
    char buffer[100];
    char hash[32];

    size_t time_str_size;

    time_rounded = round_time(time(NULL), interval);
    time_str_size = get_time_string(time_rounded, buffer, sizeof(buffer));
    if (!time_str_size) {
        assert(0 && "Buffer size is not enough");
    }

    size_t input_size = time_str_size + client_secret_size;
    char input[input_size];
    memcpy(input, buffer, time_str_size);
    memcpy(input + time_str_size, client_secret, client_secret_size);

    int rc = generate_hash(hash, secret, secret_size, input, input_size);
    if (rc) {
        return 0;
    }

    uint64_t truncated = 0;
    // Todo: handle endian
    memcpy(&truncated, hash, sizeof(truncated));
    return truncated % (int)pow(10, digit);
}