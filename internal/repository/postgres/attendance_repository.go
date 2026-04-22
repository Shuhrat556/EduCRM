package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AttendanceRepository implements repository.AttendanceRepository.
type AttendanceRepository struct {
	db *gorm.DB
}

var _ repository.AttendanceRepository = (*AttendanceRepository)(nil)

// NewAttendanceRepository constructs AttendanceRepository.
func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

// Create inserts an attendance row.
func (r *AttendanceRepository) Create(ctx context.Context, a *domain.Attendance) error {
	m, err := attendanceToModel(a)
	if err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		if isUniqueViolation(err) {
			return repository.ErrDuplicate
		}
		return err
	}
	return nil
}

// Update updates status and comment (and timestamps).
func (r *AttendanceRepository) Update(ctx context.Context, a *domain.Attendance) error {
	m, err := attendanceToModel(a)
	if err != nil {
		return err
	}
	res := r.db.WithContext(ctx).Model(&model.Attendance{}).Where("id = ?", a.ID).Updates(map[string]any{
		"status":    m.Status,
		"comment":   m.Comment,
		"updated_at": m.UpdatedAt,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindByID loads by primary key.
func (r *AttendanceRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Attendance, error) {
	var m model.Attendance
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return attendanceToDomain(&m)
}

// ListByStudent lists attendance for a student, optionally scoped to a teacher's groups.
func (r *AttendanceRepository) ListByStudent(ctx context.Context, studentID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Attendance, error) {
	q := r.db.WithContext(ctx).Model(&model.Attendance{}).Where("attendances.student_id = ?", studentID)
	if viewerTeacherID != nil {
		q = q.Joins(`JOIN teacher_group_subject_assignments tgsa ON tgsa.group_id = attendances.group_id AND tgsa.subject_id = attendances.subject_id AND tgsa.teacher_id = ?`, *viewerTeacherID)
	}
	var rows []model.Attendance
	if err := q.Order("attendances.lesson_date DESC, attendances.created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return attendanceRowsToDomain(rows)
}

// ListByGroup lists attendance for a group.
func (r *AttendanceRepository) ListByGroup(ctx context.Context, groupID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Attendance, error) {
	q := r.db.WithContext(ctx).Model(&model.Attendance{}).Where("attendances.group_id = ?", groupID)
	if viewerTeacherID != nil {
		q = q.Joins(`JOIN teacher_group_subject_assignments tgsa ON tgsa.group_id = attendances.group_id AND tgsa.subject_id = attendances.subject_id AND tgsa.teacher_id = ?`, *viewerTeacherID)
	}
	var rows []model.Attendance
	if err := q.Order("attendances.lesson_date DESC, attendances.created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return attendanceRowsToDomain(rows)
}

// ListByDateRange lists rows with lesson_date in [from, to] (inclusive, date parts only).
func (r *AttendanceRepository) ListByDateRange(ctx context.Context, from, to time.Time, viewerTeacherID *uuid.UUID) ([]domain.Attendance, error) {
	fromD := truncateUTCDate(from)
	toD := truncateUTCDate(to)
	q := r.db.WithContext(ctx).Model(&model.Attendance{}).
		Where("attendances.lesson_date >= ? AND attendances.lesson_date <= ?", fromD, toD)
	if viewerTeacherID != nil {
		q = q.Joins(`JOIN teacher_group_subject_assignments tgsa ON tgsa.group_id = attendances.group_id AND tgsa.subject_id = attendances.subject_id AND tgsa.teacher_id = ?`, *viewerTeacherID)
	}
	var rows []model.Attendance
	if err := q.Order("attendances.lesson_date ASC, attendances.group_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return attendanceRowsToDomain(rows)
}

func truncateUTCDate(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func attendanceToDomain(m *model.Attendance) (*domain.Attendance, error) {
	st, err := domain.ParseAttendanceStatus(m.Status)
	if err != nil {
		return nil, err
	}
	return &domain.Attendance{
		ID:                m.ID,
		StudentID:         m.StudentID,
		GroupID:           m.GroupID,
		SubjectID:         m.SubjectID,
		LessonDate:        truncateUTCDate(m.LessonDate),
		Status:            st,
		Comment:           m.Comment,
		MarkedByTeacherID: m.MarkedByTeacherID,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}, nil
}

func attendanceToModel(a *domain.Attendance) (*model.Attendance, error) {
	return &model.Attendance{
		ID:                a.ID,
		StudentID:         a.StudentID,
		GroupID:           a.GroupID,
		SubjectID:         a.SubjectID,
		LessonDate:        truncateUTCDate(a.LessonDate),
		Status:            string(a.Status),
		Comment:           a.Comment,
		MarkedByTeacherID: a.MarkedByTeacherID,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}, nil
}

func attendanceRowsToDomain(rows []model.Attendance) ([]domain.Attendance, error) {
	out := make([]domain.Attendance, 0, len(rows))
	for i := range rows {
		a, err := attendanceToDomain(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return out, nil
}
