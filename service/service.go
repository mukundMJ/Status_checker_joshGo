package service

import (
	"context"
	"net/http"
)

type StatusChecker interface {
	CheckSite(ctx context.Context, name string) (status bool, err error)
}

type HttpChecker struct {
}

func (h HttpChecker) CheckSite(ctx context.Context, link string) (bool, error) {
	var status bool
	_, err := http.Get(link)
	if err != nil {
		status = false
		return status, err
	}
	status = true
	return status, nil
}
