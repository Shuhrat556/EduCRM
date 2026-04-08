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

// GradeRepository implements repository.GradeRepository.
type GradeRepository struct {
	db *gorm.DB
}

var _ repository.GradeRepository = (*GradeRepository)(nil)

// NewGradeRepository constructs GradeRepository.
func NewGradeRepository(db *gorm.DB) *GradeRepository {
	return &GradeRepository{db: db}
}

// Create inserts a grade.
func (r *GradeRepository) Create(ctx context.Context, g *domain.Grade) error {
	m, err := gradeToModel(g)
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

// Update updates grade_value, comment, graded_at.
func (r *GradeRepository) Update(ctx context.Context, g *domain.Grade) error {
	m, err := gradeToModel(g)
	if err != nil {
		return err
	}
	res := r.db.WithContext(ctx).Model(&model.Grade{}).Where("id = ?", g.ID).Updates(map[string]any{
		"grade_value": m.GradeValue,
		"comment":     m.Comment,
		"graded_at":   m.GradedAt,
		"updated_at":  m.UpdatedAt,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete removes a grade by id.
func (r *GradeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Grade{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindByID loads by primary key.
func (r *GradeRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Grade, error) {
	var m model.Grade
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return gradeToDomain(&m)
}

// ListByStudent lists grades for a student.
func (r *GradeRepository) ListByStudent(ctx context.Context, studentID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Grade, error) {
	q := r.db.WithContext(ctx).Model(&model.Grade{}).Where("grades.student_id = ?", studentID)
	if viewerTeacherID != nil {
		q = q.Joins("JOIN groups ON groups.id = grades.group_id AND groups.teacher_id = ?", *viewerTeacherID)
	}
	var rows []model.Grade
	if err := q.Order("grades.week_start_date DESC, grades.created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return gradesToDomain(rows)
}

// ListByGroup lists grades for a group.
func (r *GradeRepository) ListByGroup(ctx context.Context, groupID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Grade, error) {
	q := r.db.WithContext(ctx).Model(&model.Grade{}).Where("grades.group_id = ?", groupID)
	if viewerTeacherID != nil {
		q = q.Joins("JOIN groups ON groups.id = grades.group_id AND groups.teacher_id = ?", *viewerTeacherID)
	}
	var rows []model.Grade
	if err := q.Order("grades.week_start_date DESC, grades.student_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return gradesToDomain(rows)
}

func gradeToDomain(m *model.Grade) (*domain.Grade, error) {
	gt, err := domain.ParseGradeType(m.GradeType)
	if err != nil {
		return nil, err
	}
	ws := truncateUTCDate(m.WeekStartDate)
	return &domain.Grade{
		ID:            m.ID,
		StudentID:     m.StudentID,
		TeacherID:     m.TeacherID,
		GroupID:       m.GroupID,
		WeekStartDate: ws,
		GradeType:     gt,
		GradeValue:    m.GradeValue,
		Comment:       m.Comment,
		GradedAt:      m.GradedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}, nil
}

func gradeToModel(g *domain.Grade) (*model.Grade, error) {
	return &model.Grade{
		ID:            g.ID,
		StudentID:     g.StudentID,
		TeacherID:     g.TeacherID,
		GroupID:       g.GroupID,
		WeekStartDate: truncateUTCDate(g.WeekStartDate),
		GradeType:     string(g.GradeType),
		GradeValue:    g.GradeValue,
		Comment:       g.Comment,
		GradedAt:      g.GradedAt,
		CreatedAt:     g.CreatedAt,
		UpdatedAt:     g.UpdatedAt,
	}, nil
}

func gradesToDomain(rows []model.Grade) ([]domain.Grade, error) {
	out := make([]domain.Grade, 0, len(rows))
	for i := range rows {
		g, err := gradeToDomain(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *g)
	}
	return out, nil
}
