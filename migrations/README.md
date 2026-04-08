# Database migrations

`000001_auth.*` defines `users` and `refresh_tokens` for authentication. `000002_user_active.*` adds `users.is_active`. `000003_teachers.*` adds `groups`, `teachers`, and `teacher_groups`. `000004_rooms.*` adds `rooms` for scheduling. `000005_groups_subjects.*` adds `subjects`, extends `groups` (subject, teacher, room, dates, fee, status), drops `teacher_groups`, and adds `student_group_memberships`. `000006_schedules.*` adds `schedules` (weekly slots per group/teacher/room). `000007_attendance.*` adds `user_teacher_links` (teacher login → `teachers` row) and `attendances`. `000008_grades.*` adds `grades` (weekly ratings per student/group/teacher/type). The API also runs GORM AutoMigrate for these tables on startup.

Place additional versioned SQL migrations here (for example with [golang-migrate/migrate](https://github.com/golang-migrate/migrate) or [pressly/goose](https://github.com/pressly/goose)).

Suggested naming:

- `000001_init_schema.up.sql`
- `000001_init_schema.down.sql`

Run migrations locally or in CI/CD:

```bash
make migrate-up        # or: go run ./cmd/migrate up
```

Set `MIGRATIONS_PATH` if the migrations folder is not `./migrations`. The API can skip GORM AutoMigrate in production with `AUTO_MIGRATE=false` after SQL migrations are applied.
