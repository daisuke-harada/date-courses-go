package openapi

import (
	"github.com/daisuke-harada/date-courses-go/internal/usecase"
)

func BuildGetCoursesResponse(output *usecase.GetCoursesOutput) (map[string]interface{}, error) {
	courses := make([]map[string]interface{}, len(output.Courses))
	for i, course := range output.Courses {
		courses[i] = map[string]interface{}{
			"id":            course.ID,
			"user_id":       course.UserID,
			"travel_mode":   course.TravelMode,
			"authority":     course.Authority,
			"created_at":    course.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			"updated_at":    course.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return map[string]interface{}{
		"courses": courses,
	}, nil
}
