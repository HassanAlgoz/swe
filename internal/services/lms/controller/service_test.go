package controller

import (
	"context"
	"testing"

	"github.com/hassanalgoz/swe/internal/services/lms/store"
	StorePort "github.com/hassanalgoz/swe/internal/services/lms/store/port"
	"github.com/hassanalgoz/swe/pkg/external/s3"
	"github.com/hassanalgoz/swe/pkg/services/adapters/notify"
	NotifyPort "github.com/hassanalgoz/swe/pkg/services/ports/notify"
	"github.com/hassanalgoz/swe/pkg/utils"
)

func TestCourseCRUD(t *testing.T) {
	// internal service: use fakes (even better than baked responses; stubs)
	notifyClient := notify.NewMock(notify.MockState{
		Notifications: []NotifyPort.Notification{},
	}, notify.MockFuncs{})

	// external service: use fakes (even better than baked responses; stubs)
	s3Client := s3.NewMock(s3.MockState{
		Files: map[string][]byte{},
		Tags:  map[string]map[string]string{},
	})

	// store is a direct infrastructure dependency, so: use the real thing
	dbname := utils.RandomString(30)
	storeClient := store.New(dbname)

	controller := New(storeClient, notifyClient, s3Client)

	// Test: Create
	course := StorePort.Course{
		Code:        "ICS202",
		Name:        "Test Course",
		Description: "This is a test course",
	}
	createCourseResult, err := controller.CreateCourse(context.Background(), course)
	if err != nil {
		t.Fatalf("failed to create course: %v", err)
	}

	// Test: Get
	_, err = controller.GetCourse(context.Background(), createCourseResult.ID)
	if err != nil {
		t.Fatalf("failed to get course: %v", err)
	}

	// Test: Update
	update := StorePort.Course{
		Code:        "ICS202",
		Name:        "Updated Test Course",
		Description: "This is an updated test course",
	}
	_, err = controller.UpdateCourse(context.Background(), createCourseResult.ID, update)
	if err != nil {
		t.Fatalf("failed to update course: %v", err)
	}

	// Test: Delete
	err = controller.DeleteCourse(context.Background(), createCourseResult.ID)
	if err != nil {
		t.Fatalf("failed to delete course: %v", err)
	}

	// verify that the course was deleted
	deletedCourse, err := controller.GetCourse(context.Background(), createCourseResult.ID)
	if deletedCourse != nil || err == nil {
		t.Fatalf("course was not deleted") // Note: this might err with no rows ???
	}
}

// func TestUpdateCourseErrorScenarios(t *testing.T) {
// 	// external service: use fakes (even better than baked responses; stubs)
// 	s3Client := s3.NewMock(s3.MockState{
// 		Files: map[string][]byte{},
// 		Tags:  map[string]map[string]string{},
// 	})

// 	// store is a direct infrastructure dependency, so: use the real thing
// 	dbname := utils.RandomString(30)
// 	storeClient := store.New(dbname)

// 	t.Run("send notification errs with codes.DeadlineExceeded", func(t *testing.T) {
// 		notifyClient := notify.NewMock(notify.MockState{
// 			Notifications: []NotifyPort.Notification{},
// 		}, notify.MockFuncs{
// 			SendNotification: func(in *NotifyPort.NotificationRequest) (*NotifyPort.NotificationsResponse, error) {
// 				return nil, status.Error(codes.DeadlineExceeded, "Deadline Exceeded")
// 			},
// 		})

// 		controller := New(storeClient, notifyClient, s3Client)

// 		id := uuid.MustParse("3cf74680-4b50-4ccd-9e8a-009272594858")
// 		updatedCourse := xstatus.Course{}
// 		err := controller.UpdateCourse(context.Background(), id, updatedCourse)
// 		assert.Equal(t, xstatus.ErrDeadlineExceeded{}, err)
// 	})

// 	t.Run("send notification errs with codes.Internal", func(t *testing.T) {
// 		notifyClient := notify.NewMock(notify.MockState{
// 			Notifications: []NotifyPort.Notification{},
// 		}, notify.MockFuncs{
// 			SendNotification: func(in *NotifyPort.NotificationRequest) (*NotifyPort.NotificationsResponse, error) {
// 				return nil, status.Error(codes.Internal, "Internal")
// 			},
// 		})

// 		controller := New(storeClient, notifyClient, s3Client)

// 		id := uuid.MustParse("3cf74680-4b50-4ccd-9e8a-009272594858")
// 		updatedCourse := xstatus.Course{}
// 		err := controller.UpdateCourse(context.Background(), id, updatedCourse)
// 		assert.Equal(t, xstatus.ErrInternal{}, err)
// 	})
// }

// func TestUpdateCourseUnhappy(t *testing.T) {
// 	namespace := "lms_TestUpdateCourseUnhappy"

// 	stor := store.Get(namespace)   // direct infra: the real thing
// 	notif := notify.Get(namespace) // internal service: Use fakes for internal services (even better than baked responses; stubs)
// 	controller := New(stor, notif)

// 	// TODO: test the unhappy paths
// 	courseId := uuid.MustParse("1d3d3d29-6d2c-4b5d-b3e3-1a7e0f8c4a5e")
// 	testCases := []struct {
// 		casename    string
// 		current     entities.Course
// 		update      entities.Course
// 		expectedErr error
// 	}{
// 		{
// 			casename: "successful course update",
// 			current: entities.Course{
// 				ID:          courseId,
// 				Name:        "ICS209",
// 				Description: "A beginner's course on Go programming language",
// 			},
// 			update: entities.Course{
// 				Name:        "ICS444",
// 				Description: "A beginner's course on Go programming language",
// 			},
// 			expectedErr: nil,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		// Given
// 		// TODO: rollback store state?
// 		// ...

// 		// When
// 		ctx := context.Background()
// 		err := controller.UpdateCourse(ctx, tc.current.ID, tc.update)

// 		// Then
// 		assert.Equal(t, tc.expectedErr, err)
// 	}
// }
