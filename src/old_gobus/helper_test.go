package gobus

import (
	"testing"
)

type Service struct {
}

func (s *Service) WithReply(arg int, reply *int) error {
	return nil
}

func (s *Service) PtrWithReply(arg *int, reply *int) error {
	return nil
}

func (s *Service) Batch(args []int) error {
	return nil
}

func (s *Service) ReplyNotPtr(arg int, reply int) error {
	return nil
}

func (s *Service) NoError(arg int ,reply *int) {
}

func (s *Service) BatchNotArray(arg int) error {
	return nil
}

func (s *Service) BatchNoError(args []int) {
}

func TestServiceMethod(t *testing.T) {
	doMap, batchMap := getMethods(&Service{})
	if len(doMap) != 2 {
		t.Error("do map should have 2 method")
	}
	if len(batchMap) != 1 {
		t.Error("batch map should have 1 method")
	}

	if _, ok := doMap["WithReply"]; !ok {
		t.Error("do map should have WithReply")
	}
	if _, ok := doMap["PtrWithReply"]; !ok {
		t.Error("do map should have PtrWithReply")
	}
	if _, ok := batchMap["Batch"]; !ok {
		t.Error("batch map should have Batch")
	}
}
