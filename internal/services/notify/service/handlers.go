package service

import (
	"github.com/google/uuid"
	port "github.com/hassanalgoz/swe/internal/services/notify/store/port"
)

func (s *service) GetStudentById(id uuid.UUID) (*port.Student, error) {
	result, err := s.store.GetStudentById(s.ctx, id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
