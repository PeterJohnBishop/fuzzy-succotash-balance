# fuzzy-succotash-tui

Just a little test to see if calculating,

- the 50000th number of the Fibonacci sequence
- Pi to 50000 iterations
- the number of Primes in 50000

is faster calling sequentally or concurrently.

The result was surprising:

[GIN] 2025/04/22 - 09:46:06 | 200 |   29.021334ms |             ::1 | GET      "/calc/all/inline"
[GIN] 2025/04/22 - 09:46:11 | 200 |     32.7945ms |             ::1 | GET      "/calc/all/concurrent"

until considering the cost to setup the Go routines and wait group!