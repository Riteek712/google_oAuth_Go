package server

import (
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.GET("/logout/:provider", s.oAuthLogout)
	r.GET("/auth/:provider", s.oAuthToken)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) getOAuthCallbackFunction(c *gin.Context) {
	provider := c.Param("provider")

	// Create a new request with the provider in context
	r := c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
	w := c.Writer

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error completing user authentication: %v", err)
		return
	}
	fmt.Println(user)

	c.Redirect(http.StatusFound, "http://localhost:8080")
}

func (s *Server) oAuthLogout(c *gin.Context) {
	provider := c.Param("provider")

	// Create a new request with the provider in context
	r := c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
	w := c.Writer

	gothic.Logout(w, r)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (s *Server) oAuthToken(c *gin.Context) {
	provider := c.Param("provider")

	// Create a new request with the provider in context
	r := c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
	w := c.Writer

	// Try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(w, gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

const userTemplate = `<html><body>{{.Name}}</body></html>`
