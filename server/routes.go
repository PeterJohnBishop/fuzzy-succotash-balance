package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine, port string) {

	r.GET("/", func(c *gin.Context) {
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
}
