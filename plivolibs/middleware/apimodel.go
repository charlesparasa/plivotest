package middleware

import "net/http"

type route struct {
	Name                   string
	Method                 string
	Pattern                string
	ResourcesPermissionMap map[string]uint8
	HandlerFunc            http.HandlerFunc
}

// Routes - list of route
type Routes []route

