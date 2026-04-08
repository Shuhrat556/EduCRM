package postgres

import (
	"context"
	"errors"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScheduleRepository implements repository.ScheduleRepository.
type ScheduleRepository struct {
	db *gorm.DB
}

var _ repository.ScheduleRepository = (*ScheduleRepository)(nil)

// NewScheduleRepository constructs ScheduleRepository.
func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// Create inserts a schedule row.
func (r *ScheduleRepository) Create(ctx context.Context, s *domain.Schedule) error {
	m := scheduleToModel(s)
	return r.db.WithContext(ctx).Create(m).Error
}

// Update updates all scalar fields.
func (r *ScheduleRepository) Update(ctx context.Context, s *domain.Schedule) error {
	m := scheduleToModel(s)
	return r.db.WithContext(ctx).Model(&model.Schedule{}).Where("id = ?", s.ID).Updates(map[string]any{
		"group_id":        m.GroupID,
		"teacher_id":      m.TeacherID,
		"room_id":         m.RoomID,
		"weekday":         m.Weekday,
		"start_minutes":   m.StartMinutes,
		"end_minutes":     m.EndMinutes,
		"updated_at":      m.UpdatedAt,
	}).Error
}

// Delete removes a schedule by id.
func (r *ScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Schedule{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindByID loads a schedule by primary key.
func (r *ScheduleRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Schedule, error) {
	var m model.Schedule
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return scheduleToDomain(&m), nil
}

// ListByGroup returns schedules for a group ordered by weekday and start time.
func (r *ScheduleRepository) ListByGroup(ctx context.Context, groupID uuid.UUID) ([]domain.Schedule, error) {
	var rows []model.Schedule
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).
		Order("weekday ASC, start_minutes ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return schedulesToDomain(rows)
}

// ListByTeacher returns schedules for a teacher.
func (r *ScheduleRepository) ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]domain.Schedule, error) {
	var rows []model.Schedule
	if err := r.db.WithContext(ctx).Where("teacher_id = ?", teacherID).
		Order("weekday ASC, start_minutes ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return schedulesToDomain(rows)
}

// ListByRoom returns schedules for a room.
func (r *ScheduleRepository) ListByRoom(ctx context.Context, roomID uuid.UUID) ([]domain.Schedule, error) {
	var rows []model.Schedule
	if err := r.db.WithContext(ctx).Where("room_id = ?", roomID).
		Order("weekday ASC, start_minutes ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return schedulesToDomain(rows)
}

// CountRoomOverlaps counts conflicting rows (same room & weekday, overlapping time ranges).
// Overlap: [startMin, endMin) overlaps existing [s,e) iff startMin < e && endMin > s (half-open intervals).
func (r *ScheduleRepository) CountRoomOverlaps(ctx context.Context, roomID uuid.UUID, weekday domain.Weekday, startMin, endMin int, excludeID *uuid.UUID) (int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Schedule{}).
		Where("room_id = ? AND weekday = ?", roomID, int16(weekday)).
		Where("start_minutes < ? AND end_minutes > ?", endMin, startMin)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var n int64
	err := q.Count(&n).Error
	return n, err
}

// CountTeacherOverlaps counts conflicting rows for the same teacher.
func (r *ScheduleRepository) CountTeacherOverlaps(ctx context.Context, teacherID uuid.UUID, weekday domain.Weekday, startMin, endMin int, excludeID *uuid.UUID) (int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Schedule{}).
		Where("teacher_id = ? AND weekday = ?", teacherID, int16(weekday)).
		Where("start_minutes < ? AND end_minutes > ?", endMin, startMin)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var n int64
	err := q.Count(&n).Error
	return n, err
}

func scheduleToDomain(m *model.Schedule) *domain.Schedule {
	return &domain.Schedule{
		ID:           m.ID,
		GroupID:      m.GroupID,
		TeacherID:    m.TeacherID,
		RoomID:       m.RoomID,
		Weekday:      domain.Weekday(m.Weekday),
		StartMinutes: m.StartMinutes,
		EndMinutes:   m.EndMinutes,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func scheduleToModel(s *domain.Schedule) *model.Schedule {
	return &model.Schedule{
		ID:           s.ID,
		GroupID:      s.GroupID,
		TeacherID:    s.TeacherID,
		RoomID:       s.RoomID,
		Weekday:      int16(s.Weekday),
		StartMinutes: s.StartMinutes,
		EndMinutes:   s.EndMinutes,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

func schedulesToDomain(rows []model.Schedule) ([]domain.Schedule, error) {
	out := make([]domain.Schedule, 0, len(rows))
	for i := range rows {
		out = append(out, *scheduleToDomain(&rows[i]))
	}
	return out, nil
}
