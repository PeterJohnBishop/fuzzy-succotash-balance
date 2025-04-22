package server

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine, port string) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": fmt.Sprintf("Drinking Gin on %s", port),
		})
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204) // No Content
	})

	r.GET("/apple-touch-icon.png", func(c *gin.Context) {
		c.Status(204)
	})

	r.GET("/apple-touch-icon-precomposed.png", func(c *gin.Context) {
		c.Status(204)
	})
}

func loadTestingRoutes(r *gin.Engine) {

	r.GET("/calc/fibonacci", func(c *gin.Context) {
		start := time.Now()

		nth := 50000
		result := computeFibonacci(nth) // Calculate to the nth value

		duration := time.Since(start)

		c.JSON(http.StatusOK, gin.H{
			"duration": duration.String(),
			"result":   fmt.Sprintf("%s...", result.String()[:10]), // Return just a snippet
		})
	})

	r.GET("/calc/pi", func(c *gin.Context) {
		start := time.Now()

		iterations := 50000
		result := calculatePi(iterations) // calculate pi to x iterations

		duration := time.Since(start)

		c.JSON(http.StatusOK, gin.H{
			"duration": duration.String(),
			"result":   fmt.Sprintf("Pi approximation at %d iterations: %.15f\n", iterations, result), // Return just a snippet
		})
	})

	r.GET("/calc/prime", func(c *gin.Context) {
		start := time.Now()

		limit := 50000
		result := countPrime(limit) // calculate number of prime numbers in limit

		duration := time.Since(start)

		c.JSON(http.StatusOK, gin.H{
			"duration": duration.String(),
			"result":   fmt.Sprintf("%d primes", result), // Return count
		})
	})

	r.GET("/calc/all/inline", func(c *gin.Context) {
		start := time.Now()

		nth := 50000
		sequence := computeFibonacci(nth)

		limit := 50000
		result := countPrime(limit)

		iterations := 50000
		pi := calculatePi(iterations)

		duration := time.Since(start)

		c.JSON(http.StatusOK, gin.H{
			"duration":  duration.String(),
			"Fibonacci": fmt.Sprintf("%s...", sequence.String()[:10]),
			"Pi":        fmt.Sprintf("Pi approximation at %d iterations: %.15f\n", iterations, pi),
			"Primes":    fmt.Sprintf("%d primes", result),
		})
	})

	type allResponse struct {
		Fibonacci string `json:"fibonacci"`
		Pi        string `json:"pi"`
		Primes    string `json:"primes"`
	}

	r.GET("/calc/all/concurrent", func(c *gin.Context) {
		start := time.Now()
		runtime.GOMAXPROCS(runtime.NumCPU())
		var response allResponse
		var wg sync.WaitGroup
		wg.Add(3)

		nth := 50000
		go func() {
			defer wg.Done()
			result := computeFibonacci(nth)
			response.Fibonacci = fmt.Sprintf("%s...", result.String()[:10])
		}()

		limit := 50000
		go func() {
			defer wg.Done()
			result := countPrime(limit)
			response.Primes = fmt.Sprintf("%d primes", result)
		}()

		iterations := 50000
		go func() {
			defer wg.Done()
			result := calculatePi(iterations)
			response.Pi = fmt.Sprintf("Pi approximation at %d iterations: %.15f", iterations, result)
		}()

		wg.Wait()
		duration := time.Since(start)

		c.JSON(http.StatusOK, gin.H{
			"duration": duration.String(),
			"results":  response,
		})
	})
}
