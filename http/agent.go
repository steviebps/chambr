package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/steviebps/realm/utils"
)

type Operation string

const (
	PutOperation    Operation = "put"
	GetOperation    Operation = "get"
	DeleteOperation Operation = "delete"
	ListOperation   Operation = "list"
)

type AgentRequest struct {
	*http.Request
	ID        string
	Operation Operation
	Path      string
}

func buildAgentRequest(req *http.Request) *AgentRequest {
	p := strings.TrimPrefix(req.URL.Path, "/v1/chambers")
	var op Operation

	switch req.Method {
	case http.MethodGet:
		op = GetOperation
		listStr := req.URL.Query().Get("list")
		if listStr != "" {
			list, _ := strconv.ParseBool(listStr)
			if list {
				op = ListOperation
			}
		}
	case http.MethodPost:
		op = PutOperation
	case http.MethodDelete:
		op = DeleteOperation
	case "LIST":
		op = ListOperation
	}

	return &AgentRequest{
		Request:   req,
		ID:        uuid.New().String(),
		Operation: op,
		Path:      utils.EnsureTrailingSlash(p),
	}
}
