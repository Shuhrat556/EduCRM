package domain

import "fmt"

// NotificationType categorizes in-app notifications and outbound routing hints.
type NotificationType string

const (
	NotificationTypePaymentReminder  NotificationType = "payment_reminder"
	NotificationTypePaymentReceived  NotificationType = "payment_received"
	NotificationTypeGradePosted      NotificationType = "grade_posted"
	NotificationTypeAttendanceMarked NotificationType = "attendance_marked"
	NotificationTypeScheduleChanged  NotificationType = "schedule_changed"
	NotificationTypeDocumentShared   NotificationType = "document_shared"
	NotificationTypeSystem           NotificationType = "system"
	NotificationTypeAnnouncement     NotificationType = "announcement"
	NotificationTypeMessage          NotificationType = "message"
)

var validNotificationTypes = map[NotificationType]struct{}{
	NotificationTypePaymentReminder:  {},
	NotificationTypePaymentReceived:  {},
	NotificationTypeGradePosted:      {},
	NotificationTypeAttendanceMarked: {},
	NotificationTypeScheduleChanged:  {},
	NotificationTypeDocumentShared:   {},
	NotificationTypeSystem:           {},
	NotificationTypeAnnouncement:     {},
	NotificationTypeMessage:          {},
}

// ParseNotificationType validates s.
func ParseNotificationType(s string) (NotificationType, error) {
	t := NotificationType(s)
	if _, ok := validNotificationTypes[t]; !ok {
		return "", fmt.Errorf("%w: %q", ErrInvalidNotificationType, s)
	}
	return t, nil
}
