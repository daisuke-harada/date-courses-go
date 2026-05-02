package openapi

import "github.com/daisuke-harada/date-courses-go/internal/domain/model"

// NewCoursesResponse は []*model.Course から []CourseResponseDataBody を構築します。
func NewCoursesResponse(courses []*model.Course) ([]CourseResponseData, error) {
	responses := make([]CourseResponseData, 0, len(courses))
	for _, c := range courses {
		cr, err := buildCourseResponseBody(c)
		if err != nil {
			return nil, err
		}
		responses = append(responses, cr)
	}
	return responses, nil
}

// NewCourseResponse は *model.Course から CourseResponseData を構築します。
func NewCourseResponse(course *model.Course) (CourseResponseData, error) {
	return buildCourseResponseBody(course)
}

// NewCreateCourseResponse は CourseID から CourseFormResponseData を構築します。
func NewCreateCourseResponse(courseID uint) CourseFormResponseData {
	return CourseFormResponseData{CourseId: int(courseID)}
}
