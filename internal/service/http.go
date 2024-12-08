package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RacoonMediaServer/rms-mirror/internal/config"
	"go-micro.dev/v4/logger"
)

func extractLink(query url.Values) (*url.URL, error) {
	href := query.Get("href")
	if href == "" {
		return nil, errors.New("href must be set")
	}

	decoded, err := base64.URLEncoding.DecodeString(href)
	if err != nil {
		return nil, fmt.Errorf("decode base64 encoded URL failed: %w", err)
	}

	return url.Parse(string(decoded))
}

func (s *Service) getDomainConfig(u *url.URL) (config.Domain, error) {
	domain, ok := s.cfg.AllowedDomains[u.Hostname()]
	if !ok {
		return config.Domain{}, fmt.Errorf("denied domain: %s", u.Hostname())
	}
	return domain, nil
}

func (s *Service) proxyFunc(w http.ResponseWriter, req *http.Request) {
	link, err := extractLink(req.URL.Query())
	if err != nil {
		logger.Errorf("Extract content URL failed: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	domain, err := s.getDomainConfig(link)
	if err != nil {
		logger.Errorf("URL %s is prohibited: %s", link.String(), err)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	logger.Info(domain.ContentType)
	logger.Infof("%+v", req.Header)

}
