package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
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

func requestContent(link *url.URL, userAgent, accept string, maxBytes int64) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, link.String(), nil)
	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", accept)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return &http.Response{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.ContentLength > maxBytes {
		_ = resp.Body.Close()
		return &http.Response{}, fmt.Errorf("limit exceede: %d > %d", resp.ContentLength, maxBytes)
	}

	return resp, nil
}

func (s *Service) proxyFunc(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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

	resp, err := requestContent(link, req.UserAgent(), domain.MakeAcceptHeader(), domain.LimitBytes())
	if err != nil {
		logger.Errorf("Fetch %s failed: %s", link.String(), err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Add("Content-Type", resp.Header.Get("Content-Type"))
	_, _ = io.CopyN(w, resp.Body, domain.LimitBytes())
}
