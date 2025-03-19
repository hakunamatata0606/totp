#pragma once

#include <stddef.h>
#include <stdint.h>

uint64_t generate_otp(const char* secret, size_t secret_size, const char* client_secret, size_t client_secret_size, uint64_t interval, uint8_t digit);
