package proxy

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ProxyHandler creates a reverse proxy handler for a backend service
func ProxyHandler(targetURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build target URL
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		// Construct full target URL
		fullURL := targetURL + path
		if rawQuery != "" {
			fullURL += "?" + rawQuery
		}

		// Debug log: where the request is being proxied
		log.Printf("Proxying %s %s -> %s", c.Request.Method, c.Request.URL.String(), fullURL)

		// Create new request to backend service
		req, err := http.NewRequest(c.Request.Method, fullURL, c.Request.Body)
		if err != nil {
			log.Printf("Error creating proxy request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create proxy request"})
			return
		}

		// Copy headers from original request
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error forwarding request to %s: %v", targetURL, err)
			c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}

		// Set status code
		c.Status(resp.StatusCode)

		// Copy response body
		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			log.Printf("Error copying response body: %v", err)
		}
	}
}

// StripPrefixProxy strips the prefix from the path before forwarding
func StripPrefixProxy(prefix, targetURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Strip prefix from path
		path := c.Request.URL.Path
		if strings.HasPrefix(path, prefix) {
			path = strings.TrimPrefix(path, prefix)
		}

		rawQuery := c.Request.URL.RawQuery

		// Construct full target URL
		fullURL := targetURL + path
		if rawQuery != "" {
			fullURL += "?" + rawQuery
		}

		// Create new request to backend service
		req, err := http.NewRequest(c.Request.Method, fullURL, c.Request.Body)
		if err != nil {
			log.Printf("Error creating proxy request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create proxy request"})
			return
		}

		// Copy headers from original request
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error forwarding request to %s: %v", targetURL, err)
			c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}

		// Set status code
		c.Status(resp.StatusCode)

		// Copy response body
		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			log.Printf("Error copying response body: %v", err)
		}
	}
}
