package store

import (
	"fmt"
	"jira_crawler/cli_shell/crawler"
	"jira_crawler/cli_shell/model"
	"sort"
)

type MemoryStore struct {
	crawlerEngines map[string]crawler.Engine
	sessions       map[string]model.CrawlSession
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		crawlerEngines: make(map[string]crawler.Engine),
		sessions:       make(map[string]model.CrawlSession),
	}
}

func (s *MemoryStore) RegisterCrawlerEngine(engine crawler.Engine) error {
	if engine == nil {
		return fmt.Errorf("crawler engine cannot be nil")
	}

	definition := engine.Definition()

	if definition.ID == "" {
		return fmt.Errorf("crawler definition id is required")
	}

	if _, exists := s.crawlerEngines[definition.ID]; exists {
		return fmt.Errorf("crawler definition %q already registered", definition.ID)
	}

	s.crawlerEngines[definition.ID] = engine
	return nil
}

func (s *MemoryStore) GetCrawlerEngine(id string) (crawler.Engine, error) {
	engine, exists := s.crawlerEngines[id]
	if !exists {
		return nil, fmt.Errorf("crawler definition %q not found", id)
	}

	return engine, nil
}

func (s *MemoryStore) ListCrawlerDefinitions() []model.CrawlerDefinition {
	result := make([]model.CrawlerDefinition, 0, len(s.crawlerEngines))

	for _, engine := range s.crawlerEngines {
		result = append(result, engine.Definition())
	}

	sort.Slice(
		result, func(i, j int) bool {
			return result[i].ID < result[j].ID
		},
	)

	return result
}

func (s *MemoryStore) AddSession(session model.CrawlSession) error {
	if session.Name == "" {
		return fmt.Errorf("session name is required")
	}

	if _, exists := s.sessions[session.Name]; exists {
		return fmt.Errorf("session %q already exists", session.Name)
	}

	s.sessions[session.Name] = session
	return nil
}

func (s *MemoryStore) UpdateSession(session model.CrawlSession) error {
	if session.Name == "" {
		return fmt.Errorf("session name is required")
	}

	if _, exists := s.sessions[session.Name]; !exists {
		return fmt.Errorf("session %q not found", session.Name)
	}

	s.sessions[session.Name] = session
	return nil
}

func (s *MemoryStore) RemoveSession(name string) error {
	if _, exists := s.sessions[name]; !exists {
		return fmt.Errorf("session %q not found", name)
	}

	delete(s.sessions, name)
	return nil
}

func (s *MemoryStore) GetSession(name string) (model.CrawlSession, error) {
	session, exists := s.sessions[name]
	if !exists {
		return model.CrawlSession{}, fmt.Errorf("session %q not found", name)
	}

	return session, nil
}

func (s *MemoryStore) ListSessions() []model.CrawlSession {
	result := make([]model.CrawlSession, 0, len(s.sessions))

	for _, session := range s.sessions {
		result = append(result, session)
	}

	sort.Slice(
		result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		},
	)

	return result
}
