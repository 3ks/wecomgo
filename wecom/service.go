package wecom

import (
	"context"
	"net/http"
)

type service struct {
	client *client
	ctx    context.Context
}

func (s *service) doRequest(req *http.Request, result iBaseResponse) (err error) {
	if s.ctx != nil {
		req = req.WithContext(s.ctx)
	}
	err = s.client.do(req, result)
	if err != nil {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
		}
		return err
	}
	return nil
}
