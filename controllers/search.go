package controllers

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/quorumsco/contacts/models"
	"github.com/quorumsco/logs"

	"gopkg.in/olivere/elastic.v2"
)

type Search struct {
	Client *elastic.Client
}

func (s *Search) Index(args models.ContactArgs, reply *models.ContactReply) error {
	id := strconv.Itoa(int(args.Contact.ID))
	if id == "" {
		logs.Error("id is nil")
		return errors.New("id is nil")
	}

	_, err := s.Client.Index().
		Index("contacts").
		Type("contact").
		Id(id).
		BodyJson(reply.Contact).
		Do()
	if err != nil {
		logs.Critical(err)
		return err
	}

	return nil
}

func (s *Search) UnIndex(args models.ContactArgs, reply *models.ContactReply) error {
	id := strconv.Itoa(int(args.Contact.ID))
	if id == "" {
		logs.Error("id is nil")
		return errors.New("id is nil")
	}

	_, err := s.Client.Delete().
		Index("contacts").
		Type("contact").
		Id(id).
		Do()
	if err != nil {
		logs.Critical(err)
		return err
	}

	return nil
}

func (s *Search) SearchContacts(args models.SearchArgs, reply *models.SearchReply) error {
	termQuery := elastic.NewMultiMatchQuery(args.Search.Query, args.Search.Field, "firstname")
	termQuery = termQuery.Type("cross_fields")
	termQuery = termQuery.Operator("and")
	searchResult, err := s.Client.Search().
		Index("contacts").
		Query(&termQuery).
		Sort("surname", true).
		Pretty(true).
		Do()
	if err != nil {
		logs.Critical(err)
		return err
	}

	if searchResult.Hits != nil {
		for _, hit := range searchResult.Hits.Hits {
			var c models.Contact
			err := json.Unmarshal(*hit.Source, &c)
			if err != nil {
				logs.Error(err)
				return err
			}
			reply.Contacts = append(reply.Contacts, c)
		}
	} else {
		reply.Contacts = nil
	}

	return nil
}
