# Programs for generating plots of Legendre function of the first kind, with arbitrary (non-integer) order.

## Example - `plot-legendre`: generates plots for $P_\nu(x)$ with non-integer order $\nu$
```
./plot-legendre (main) ▶ go run main.go --legendre-err-bound 1e-5
/var/folders/_0/2d8v_l8x5r947l5f35hdx0yw0000gq/T/plot-legendre.m
./plot-legendre (main) ▶ octave --persist /var/folders/_0/2d8v_l8x5r947l5f35hdx0yw0000gq/T/plot-legendre.m
```
will generate plots for $P_\nu(x)$ with $\nu$=[0.05, 0.25, 0.75, 1, 1.25, 3, 3.75, 4].

<img width="1354" alt="Screenshot 2023-05-15 at 16 11 43" src="https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/e52dc246-8b89-4b11-a0ef-5602de387fa5">

## Example - `search-zero`: does binary search for legendre function zeros, and reproduces Jackson figure 3-6.
```
./search-zero (main*) ▶ go run main.go --legendre-err-bound=1e-10 --root-err-bound=1e-6 --prec=1000
𝜈=1.01, 19 iterations (converged by value)
𝜈=1.02, 18 iterations (converged by value)
𝜈=1.03, 19 iterations (converged by value)
𝜈=1.04, 18 iterations (converged by value)
𝜈=1.05, 19 iterations (converged by value)
𝜈=1.06, 17 iterations (converged by value)
𝜈=1.07, 18 iterations (converged by value)
𝜈=1.08, 18 iterations (converged by value)
𝜈=1.09, 19 iterations (converged by value)
𝜈=1.10, 18 iterations (converged by value)
...
𝜈=0.18, 15 iterations (converged by range)
𝜈=0.17, 15 iterations (converged by range)
𝜈=0.16, 14 iterations (converged by range)
𝜈=0.15, 14 iterations (converged by range)
𝜈=0.14, 13 iterations (converged by range)
𝜈=0.13, 12 iterations (converged by range)
𝜈=0.12, 11 iterations (converged by range)
𝜈=0.11, 10 iterations (converged by range)
𝜈=0.10, 9 iterations (converged by range)
𝜈=0.09, 8 iterations (converged by range)
𝜈=0.08, 6 iterations (converged by range)
𝜈=0.07, 4 iterations (converged by range)
𝜈=0.06, 2 iterations (converged by range)
𝜈=0.05, 1 iterations (converged by range)
/var/folders/_0/2d8v_l8x5r947l5f35hdx0yw0000gq/T/legendre-zeros.m
./search-zero (main) ▶ octave --persist /var/folders/_0/2d8v_l8x5r947l5f35hdx0yw0000gq/T/legendre-zeros.m
```
will plot the solution for $P_{\nu}(\cos\beta)=0$, and reproduce Figure 3-6.

<img width="1161" alt="Screenshot 2023-05-15 at 16 16 01" src="https://github.com/euphoricrhino/jackson-em-notes/assets/107862003/8cceada5-c58e-4891-9b26-131ec8e60d67">



