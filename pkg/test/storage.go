package test

import (
	"errors"

	"github.com/jhillyerd/inbucket/pkg/storage"
)

// StoreStub stubs storage.Store for testing.
type StoreStub struct {
	storage.Store
	mailboxes map[string][]storage.StoreMessage
	deleted   map[storage.StoreMessage]struct{}
}

// NewStore creates a new StoreStub.
func NewStore() *StoreStub {
	return &StoreStub{
		mailboxes: make(map[string][]storage.StoreMessage),
		deleted:   make(map[storage.StoreMessage]struct{}),
	}
}

// AddMessage adds a message to the specified mailbox.
func (s *StoreStub) AddMessage(m storage.StoreMessage) (id string, err error) {
	mb := m.Mailbox()
	msgs := s.mailboxes[mb]
	s.mailboxes[mb] = append(msgs, m)
	return m.ID(), nil
}

// GetMessage gets a message by ID from the specified mailbox.
func (s *StoreStub) GetMessage(mailbox, id string) (storage.StoreMessage, error) {
	if mailbox == "messageerr" {
		return nil, errors.New("internal error")
	}
	for _, m := range s.mailboxes[mailbox] {
		if m.ID() == id {
			return m, nil
		}
	}
	return nil, storage.ErrNotExist
}

// GetMessages gets all the messages for the specified mailbox.
func (s *StoreStub) GetMessages(mailbox string) ([]storage.StoreMessage, error) {
	if mailbox == "messageserr" {
		return nil, errors.New("internal error")
	}
	return s.mailboxes[mailbox], nil
}

// RemoveMessage deletes a message by ID from the specified mailbox.
func (s *StoreStub) RemoveMessage(mailbox, id string) error {
	mb, ok := s.mailboxes[mailbox]
	if ok {
		var msg storage.StoreMessage
		for i, m := range mb {
			if m.ID() == id {
				msg = m
				s.mailboxes[mailbox] = append(mb[:i], mb[i+1:]...)
				break
			}
		}
		if msg != nil {
			s.deleted[msg] = struct{}{}
			return nil
		}
	}
	return storage.ErrNotExist
}

// VisitMailboxes accepts a function that will be called with the messages in each mailbox while it
// continues to return true.
func (s *StoreStub) VisitMailboxes(f func([]storage.StoreMessage) (cont bool)) error {
	for _, v := range s.mailboxes {
		if !f(v) {
			return nil
		}
	}
	return nil
}

// MessageDeleted returns true if the specified message was deleted
func (s *StoreStub) MessageDeleted(m storage.StoreMessage) bool {
	_, ok := s.deleted[m]
	return ok
}
