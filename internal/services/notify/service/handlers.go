package service

import (
	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/micro/notify/service/store"
)

func (s *service) GetStudentById(id uuid.UUID) (*store.Student, error) {
	result, err := s.store.GetStudentById(s.ctx, id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
