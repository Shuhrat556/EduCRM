package dto

// RoomListQuery binds GET /rooms query parameters.
type RoomListQuery struct {
	PaginationQuery
	Status string `form:"status" example:"active" binding:"omitempty,oneof=active inactive"`
	Q      string `form:"q" example:"Lab A"`
}

// UserListQuery binds GET /users query parameters.
type UserListQuery struct {
	PaginationQuery
	Role string `form:"role" example:"teacher"`
	Q    string `form:"q" example:"@school.edu"`
}

// TeacherListQuery binds GET /teachers query parameters.
type TeacherListQuery struct {
	PaginationQuery
	Status string `form:"status" example:"active" binding:"omitempty,oneof=active inactive"`
	Q      string `form:"q" example:"Jane"`
}

// GroupListQuery binds GET /groups query parameters.
type GroupListQuery struct {
	PaginationQuery
	Q         string `form:"q" example:"Math"`
	Status    string `form:"status" example:"active" binding:"omitempty,oneof=active inactive"`
	TeacherID string `form:"teacher_id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"omitempty,uuid"`
	SubjectID string `form:"subject_id" binding:"omitempty,uuid"`
	RoomID    string `form:"room_id" binding:"omitempty,uuid"`
}

// ScheduleListQuery binds GET /schedules (exactly one of group_id, teacher_id, room_id — validated in handler).
type ScheduleListQuery struct {
	GroupID   string `form:"group_id" binding:"omitempty,uuid"`
	TeacherID string `form:"teacher_id" binding:"omitempty,uuid"`
	RoomID    string `form:"room_id" binding:"omitempty,uuid"`
}

// PaymentListQuery binds GET /payments (staff).
type PaymentListQuery struct {
	PaginationQuery
	Q           string `form:"q"`
	StudentID   string `form:"student_id" binding:"omitempty,uuid"`
	GroupID     string `form:"group_id" binding:"omitempty,uuid"`
	MonthFor    string `form:"month_for" example:"2025-04"`
	Status      string `form:"status"`
	PaymentType string `form:"payment_type"`
	IsFree      string `form:"is_free" example:"false"`
}

// PaymentHistoryQuery binds GET /payments/history.
type PaymentHistoryQuery struct {
	PaginationQuery
	StudentID string `form:"student_id" binding:"omitempty,uuid"`
}

// FileListQuery binds GET /files.
type FileListQuery struct {
	PaginationQuery
	OwnerType string `form:"owner_type" example:"student_photo" binding:"required,oneof=student_photo teacher_photo document"`
	OwnerID   string `form:"owner_id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required,uuid"`
}

// NotificationListQuery binds GET /notifications.
type NotificationListQuery struct {
	PaginationQuery
	UserID     string `form:"user_id" binding:"omitempty,uuid"`
	UnreadOnly string `form:"unread_only" example:"false"`
}

// GradeListQuery binds GET /grades.
type GradeListQuery struct {
	StudentID string `form:"student_id" binding:"omitempty,uuid"`
	GroupID   string `form:"group_id" binding:"omitempty,uuid"`
}

// AttendanceListQuery binds GET /attendance.
type AttendanceListQuery struct {
	StudentID string `form:"student_id" binding:"omitempty,uuid"`
	GroupID   string `form:"group_id" binding:"omitempty,uuid"`
	From      string `form:"from" example:"2025-04-01"`
	To        string `form:"to" example:"2025-04-30"`
}
