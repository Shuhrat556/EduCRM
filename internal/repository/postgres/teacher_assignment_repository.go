package postgres

import (
	"context"

	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TeacherAssignmentRepository implements repository.TeacherAssignmentRepository.
type TeacherAssignmentRepository struct {
	db *gorm.DB
}

var _ repository.TeacherAssignmentRepository = (*TeacherAssignmentRepository)(nil)

// NewTeacherAssignmentRepository constructs the repository.
func NewTeacherAssignmentRepository(db *gorm.DB) *TeacherAssignmentRepository {
	return &TeacherAssignmentRepository{db: db}
}

// Exists implements repository.TeacherAssignmentRepository.
func (r *TeacherAssignmentRepository) Exists(ctx context.Context, teacherID, groupID, subjectID uuid.UUID) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.TeacherGroupSubjectAssignment{}).
		Where("teacher_id = ? AND group_id = ? AND subject_id = ?", teacherID, groupID, subjectID).
		Count(&n).Error
	return n > 0, err
}

// HasAnyAssignmentOnGroup implements repository.TeacherAssignmentRepository.
func (r *TeacherAssignmentRepository) HasAnyAssignmentOnGroup(ctx context.Context, teacherID, groupID uuid.UUID) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.TeacherGroupSubjectAssignment{}).
		Where("teacher_id = ? AND group_id = ?", teacherID, groupID).
		Count(&n).Error
	return n > 0, err
}

// ListByTeacher implements repository.TeacherAssignmentRepository.
func (r *TeacherAssignmentRepository) ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]repository.TeacherAssignmentRow, error) {
	var rows []model.TeacherGroupSubjectAssignment
	if err := r.db.WithContext(ctx).Where("teacher_id = ?", teacherID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]repository.TeacherAssignmentRow, 0, len(rows))
	for i := range rows {
		out = append(out, repository.TeacherAssignmentRow{
			GroupID:   rows[i].GroupID,
			SubjectID: rows[i].SubjectID,
		})
	}
	return out, nil
}
