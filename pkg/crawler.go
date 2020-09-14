package app

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"log"
)

type ResponseContext struct {
	Context string
	//Graphs  []map[string]json.RawMessage
	Graphs []map[string]string
}

type Graph struct {
	Entry map[string]json.RawMessage
}

func (s *Server) GetRecipe(url string) error {
	s.crawler.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		//r := Response{}

		var smelly interface{}

		if err := json.Unmarshal([]byte(e.Text), &smelly); err != nil {
			log.Printf("this is the err one: %v", err)
		}
		log.Printf("this is the smelly: %v", smelly)

		//var tmp map[string]json.RawMessage
		//if err := json.Unmarshal([]byte(e.Text), &tmp); err != nil {
		//	log.Printf("this is the err one: %v", err)
		//}

		log.Printf("this is the tmp: %v", e.Text)
		//
		//for _, s := range tmp {
		//	var innerTmp []map[string]json.RawMessage
		//	if err := json.Unmarshal(s, &innerTmp); err != nil {
		//		log.Printf("this is the err one: %v", err)
		//	}
		//	log.Printf("this is len: %v", len(innerTmp))
		//
		//	for _, s := range innerTmp {
		//		var innerInnerTmp []map[string]json.RawMessage
		//
		//		if err := json.Unmarshal(s["recipeInstructions"], &innerInnerTmp); err != nil {
		//			log.Printf("this is the err one: %v", err)
		//		}
		//
		//		log.Printf("this is recipeIni: %v", innerInnerTmp)
		//	}
		//}
		//var r ResponseContext
		//json.Unmarshal(tmp["@graph"], &r.Graphs)
		//
		//log.Printf("this is the r: %v", r)

		//var objMap map[string]json.RawMessage
		//err := json.Unmarshal([]byte(e.Text), &objMap)
		//
		//if err != nil {
		//	log.Printf("this is the err one: %v", err)
		//}

		//var nestedDecoded []interface{}
		//json.Unmarshal(objMap["@recipeIngredient"], &nestedDecoded)
		//
		//log.Printf("this is nested: %v", nestedDecoded)
		//var decoded []interface{}
		//err = json.Unmarshal(objMap["@graph"], &decoded)
		//if err != nil {
		//	log.Printf("this is the err two: %v", err)
		//}
		//
		//for _, s := range objMap {
		//	var data map[string]string
		//	err := json.Unmarshal(s, &data)
		//
		//	if err != nil {
		//		log.Printf("this is the err from loop: %v", err)
		//	}
		//
		//	log.Printf("this is the loop value: %v", s)
		//}

		//log.Printf("this is the err for rC: %v", decoded)
	})

	err := s.crawler.Visit(url)

	if err != nil {
		return err
	}

	return nil
}
