package app

import (
	"context"
	"fmt"
	"jira_crawler/cli_shell/crawler"
	"jira_crawler/cli_shell/model"
	"jira_crawler/cli_shell/store"
	"time"
)

type Service struct {
	store *store.MemoryStore
}

func NewService(store *store.MemoryStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) RegisterCrawlerEngine(engine crawler.Engine) error {
	return s.store.RegisterCrawlerEngine(engine)
}

func (s *Service) ListCrawlerDefinitions() []model.CrawlerDefinition {
	return s.store.ListCrawlerDefinitions()
}

func (s *Service) ShowCrawlerDefinition(id string) (model.CrawlerDefinition, error) {
	engine, err := s.store.GetCrawlerEngine(id)
	if err != nil {
		return model.CrawlerDefinition{}, err
	}

	return engine.Definition(), nil
}

func (s *Service) AddSession(
	name string,
	crawlerDefinitionID string,
	url string,
) error {
	if name == "" {
		return fmt.Errorf("session name is required")
	}

	if crawlerDefinitionID == "" {
		return fmt.Errorf("crawler definition id is required")
	}

	if url == "" {
		return fmt.Errorf("session url is required")
	}

	if _, err := s.store.GetCrawlerEngine(crawlerDefinitionID); err != nil {
		return err
	}

	now := time.Now()

	session := model.CrawlSession{
		ID:                  newID("session"),
		Name:                name,
		CrawlerDefinitionID: crawlerDefinitionID,
		URL:                 url,
		Status:              model.CrawlStatusPending,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	return s.store.AddSession(session)
}

func (s *Service) UpdateSession(name string, newName string, newURL string) error {
	session, err := s.store.GetSession(name)
	if err != nil {
		return err
	}

	if newName != "" && newName != name {
		if err := s.store.RemoveSession(name); err != nil {
			return err
		}

		session.Name = newName
	}

	if newURL != "" {
		session.URL = newURL
	}

	session.UpdatedAt = time.Now()

	if newName != "" && newName != name {
		return s.store.AddSession(session)
	}

	return s.store.UpdateSession(session)
}

func (s *Service) RemoveSession(name string) error {
	return s.store.RemoveSession(name)
}

func (s *Service) ListSessions() []model.CrawlSession {
	return s.store.ListSessions()
}

func (s *Service) ShowSession(name string) (model.CrawlSession, error) {
	return s.store.GetSession(name)
}

func (s *Service) RunSession(ctx context.Context, name string) error {
	session, err := s.store.GetSession(name)
	if err != nil {
		return err
	}

	engine, err := s.store.GetCrawlerEngine(session.CrawlerDefinitionID)
	if err != nil {
		return err
	}

	now := time.Now()

	session.Status = model.CrawlStatusRunning
	session.LastRunAt = &now
	session.LastFinishedAt = nil
	session.LastError = ""
	session.UpdatedAt = now

	if err := s.store.UpdateSession(session); err != nil {
		return err
	}

	result, crawlErr := engine.Crawl(ctx, session.URL)

	finishedAt := time.Now()
	session.LastFinishedAt = &finishedAt
	session.UpdatedAt = finishedAt

	if crawlErr != nil {
		session.Status = model.CrawlStatusFailed
		session.LastError = crawlErr.Error()

		_ = s.store.UpdateSession(session)

		return crawlErr
	}

	session.Status = model.CrawlStatusSucceeded
	session.Result = result
	session.LastError = ""

	return s.store.UpdateSession(session)
}

func (s *Service) RerunSession(ctx context.Context, name string) error {
	return s.RunSession(ctx, name)
}

func (s *Service) ExportSession(
	name string,
	format string,
	outputPath string,
) error {
	session, err := s.store.GetSession(name)
	if err != nil {
		return err
	}

	if session.Status != model.CrawlStatusSucceeded {
		return fmt.Errorf(
			"session %q is not ready for export; current status is %s",
			name,
			session.Status,
		)
	}

	if session.Result == nil {
		return fmt.Errorf("session %q has no crawl result", name)
	}

	engine, err := s.store.GetCrawlerEngine(session.CrawlerDefinitionID)
	if err != nil {
		return err
	}

	switch format {
	case "csv":
		return engine.ExportCSV(outputPath, session.Result)

	case "xlsx":
		return engine.ExportXLSX(outputPath, session.Result)

	default:
		return fmt.Errorf("unsupported export format %q", format)
	}
}

func newID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
