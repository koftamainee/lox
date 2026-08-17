#ifndef CLOX_ALLOC_H_
#define CLOX_ALLOC_H_

#include <stddef.h>

typedef void *(MallocFn)(size_t bytes, void *ctx);

typedef void *(ReallocFn)(void *ptr, size_t new_size, void *ctx);

typedef void(FreeFn)(void *ptr, void *ctx);

typedef struct {
  MallocFn  *malloc;
  ReallocFn *realloc;
  FreeFn    *free;

  void *ctx;
} Allocator;

#endif // CLOX_ALLOC_H_
