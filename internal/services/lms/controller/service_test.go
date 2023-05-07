package controller

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/pkg/config"
	"github.com/hassanalgoz/swe/pkg/entities"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.SetupTestConfig()
}

func TestCourseCRUD(t *testing.T) {
	ctrl := Get("lms_TestCourseCRUD")
	// Test: Create
	course := entities.Course{
		Name:        "Test Course",
		Description: "This is a test course",
	}
	id, err := ctrl.CreateCourse(context.Background(), course)
	if err != nil {
		t.Fatalf("failed to create course: %v", err)
	}

	// Test: Get
	createdCourse, err := ctrl.GetCourseById(context.Background(), *id)
	if err != nil {
		t.Fatalf("failed to get course: %v", err)
	}

	// verify that the created course matches the original course
	if createdCourse.Name != course.Name || createdCourse.Description != course.Description {
		t.Fatalf("created course does not match original course")
	}

	// Test: Update
	updatedCourse := entities.Course{
		Name:        "Updated Test Course",
		Description: "This is an updated test course",
	}
	err = ctrl.UpdateCourse(context.Background(), *id, updatedCourse)
	if err != nil {
		t.Fatalf("failed to update course: %v", err)
	}

	// verify that the course was updated
	updatedCourseFromDB, err := ctrl.GetCourseById(context.Background(), *id)
	if err != nil {
		t.Fatalf("failed to get updated course: %v", err)
	}

	// verify that the updated course matches the updated course
	if updatedCourseFromDB.Name != updatedCourse.Name || updatedCourseFromDB.Description != updatedCourse.Description {
		t.Fatalf("updated course does not match expected course")
	}

	// Test: Delete
	err = ctrl.DeleteCourse(context.Background(), *id)
	if err != nil {
		t.Fatalf("failed to delete course: %v", err)
	}

	// verify that the course was deleted
	deletedCourse, err := ctrl.GetCourseById(context.Background(), *id)
	if deletedCourse != nil || err == nil {
		t.Fatalf("course was not deleted")
	}
}

func TestUpdateCourseUnhappy(t *testing.T) {
	ctrl := Get("lms_TestUpdateCourseUnhappy")
	// TODO: test the unhappy paths
	courseId := uuid.MustParse("1d3d3d29-6d2c-4b5d-b3e3-1a7e0f8c4a5e")
	testCases := []struct {
		casename    string
		current     entities.Course
		update      entities.Course
		expectedErr error
	}{
		{
			casename: "successful course update",
			current: entities.Course{
				ID:          courseId,
				Name:        "ICS209",
				Description: "A beginner's course on Go programming language",
			},
			update: entities.Course{
				Name:        "ICS444",
				Description: "A beginner's course on Go programming language",
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		// Given
		// TODO: rollback store state?
		// ...

		// When
		ctx := context.Background()
		err := ctrl.UpdateCourse(ctx, tc.current.ID, tc.update)

		// Then
		assert.Equal(t, tc.expectedErr, err)
	}
}
