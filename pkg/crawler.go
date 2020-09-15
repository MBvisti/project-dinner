package app

import (
	"github.com/gocolly/colly/v2"
)

// TODO: for future reference - look into json.RawMessage
func (s *Server) GetRecipe(url string) error {
	s.crawler.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {

	})

	err := s.crawler.Visit(url)

	if err != nil {
		return err
	}

	return nil
}
