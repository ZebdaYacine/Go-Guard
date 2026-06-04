package api

import (
	"log"
	"strings"
)

type PublicRoute struct {
	Path    string
	Methods []string
}

func GetPublicEndpoints() []PublicRoute {
	return []PublicRoute{
		{Path: "/health", Methods: []string{"GET"}},
		{Path: "/ready", Methods: []string{"GET"}},
		{Path: "/api/guest/public/info", Methods: []string{"GET"}},
		{Path: "/api/auth/register", Methods: []string{"POST"}},
		{Path: "/api/auth/login", Methods: []string{"POST"}},
	}
}

var PublicEndpoints = map[string][]string{
	"/health":                {"GET"},
	"/ready":                 {"GET"},
	"/api/guest/public/info": {"GET"},
	"/api/auth/register":     {"POST"},
	"/api/auth/login":        {"POST"},
}

func IsPublicEndpoint(path, method string) bool {
	log.Printf("Checking if %s %s is public", method, path)
	log.Printf("Available public endpoints: %+v", PublicEndpoints)

	if methods, exists := PublicEndpoints[path]; exists {
		log.Printf("Path found in public endpoints. Allowed methods: %v", methods)
		for _, m := range methods {
			if m == method {
				log.Printf("✅ Match found for %s %s", method, path)
				return true
			}
		}
		log.Printf("❌ Path found but method %s not allowed", method)
		return true
	}
	log.Printf("❌ Path not found in public endpoints")
	return false
}

// Optional: Add case-insensitive method comparison
func IsPublicEndpointCaseInsensitive(path, method string) bool {
	// Convert method to uppercase for comparison
	method = strings.ToUpper(method)

	if methods, exists := PublicEndpoints[path]; exists {
		for _, m := range methods {
			if m == method {
				return true
			}
		}
	}
	return false
}
