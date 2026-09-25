package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/royandi/gowatch/backend/internal/model"
)

type MonitorRepository struct {
	db *pgxpool.Pool
}

func NewMonitorRepository(
	db *pgxpool.Pool,
) *MonitorRepository {
	return &MonitorRepository{
		db: db,
	}
}

func (r *MonitorRepository) Create(
	ctx context.Context,
	userID int64,
	request model.CreateMonitorRequest,
	isActive bool,
) (model.Monitor, error) {
	query := `
		INSERT INTO monitors (
			user_id,
			name,
			url,
			monitor_type,
			interval_seconds,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			user_id,
			name,
			url,
			monitor_type,
			interval_seconds,
			is_active,
			created_at,
			updated_at
	`

	var monitor model.Monitor

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		request.Name,
		request.URL,
		request.MonitorType,
		request.IntervalSeconds,
		isActive,
	).Scan(
		&monitor.ID,
		&monitor.UserID,
		&monitor.Name,
		&monitor.URL,
		&monitor.MonitorType,
		&monitor.IntervalSeconds,
		&monitor.IsActive,
		&monitor.CreatedAt,
		&monitor.UpdatedAt,
	)

	return monitor, err
}

func (r *MonitorRepository) FindAllByUser(
	ctx context.Context,
	userID int64,
) ([]model.Monitor, error) {
	query := `
		SELECT
			id,
			user_id,
			name,
			url,
			monitor_type,
			interval_seconds,
			is_active,
			created_at,
			updated_at
		FROM monitors
		WHERE user_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	monitors := make([]model.Monitor, 0)

	for rows.Next() {
		var monitor model.Monitor

		err := rows.Scan(
			&monitor.ID,
			&monitor.UserID,
			&monitor.Name,
			&monitor.URL,
			&monitor.MonitorType,
			&monitor.IntervalSeconds,
			&monitor.IsActive,
			&monitor.CreatedAt,
			&monitor.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		monitors = append(
			monitors,
			monitor,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return monitors, nil
}

func (r *MonitorRepository) FindByIDAndUser(
	ctx context.Context,
	id int64,
	userID int64,
) (model.Monitor, error) {
	query := `
		SELECT
			id,
			user_id,
			name,
			url,
			monitor_type,
			interval_seconds,
			is_active,
			created_at,
			updated_at
		FROM monitors
		WHERE id = $1
		AND user_id = $2
	`

	var monitor model.Monitor

	err := r.db.QueryRow(
		ctx,
		query,
		id,
		userID,
	).Scan(
		&monitor.ID,
		&monitor.UserID,
		&monitor.Name,
		&monitor.URL,
		&monitor.MonitorType,
		&monitor.IntervalSeconds,
		&monitor.IsActive,
		&monitor.CreatedAt,
		&monitor.UpdatedAt,
	)

	return monitor, err
}

func (r *MonitorRepository) Update(
	ctx context.Context,
	monitor model.Monitor,
) (model.Monitor, error) {
	query := `
		UPDATE monitors
		SET
			name = $1,
			url = $2,
			monitor_type = $3,
			interval_seconds = $4,
			is_active = $5,
			updated_at = NOW()
		WHERE id = $6
		AND user_id = $7
		RETURNING
			id,
			user_id,
			name,
			url,
			monitor_type,
			interval_seconds,
			is_active,
			created_at,
			updated_at
	`

	var updated model.Monitor

	err := r.db.QueryRow(
		ctx,
		query,
		monitor.Name,
		monitor.URL,
		monitor.MonitorType,
		monitor.IntervalSeconds,
		monitor.IsActive,
		monitor.ID,
		monitor.UserID,
	).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.Name,
		&updated.URL,
		&updated.MonitorType,
		&updated.IntervalSeconds,
		&updated.IsActive,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	return updated, err
}

func (r *MonitorRepository) Delete(
	ctx context.Context,
	id int64,
	userID int64,
) (bool, error) {
	commandTag, err := r.db.Exec(
		ctx,
		`
			DELETE FROM monitors
			WHERE id = $1
			AND user_id = $2
		`,
		id,
		userID,
	)

	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}

func (r *MonitorRepository) FindDueForCheck(
	ctx context.Context,
) ([]model.Monitor, error) {
	query := `
		SELECT
			m.id,
			m.user_id,
			m.name,
			m.url,
			m.monitor_type,
			m.interval_seconds,
			m.is_active,
			m.created_at,
			m.updated_at
		FROM monitors m
		LEFT JOIN LATERAL (
			SELECT
				mr.checked_at
			FROM monitor_results mr
			WHERE mr.monitor_id = m.id
			ORDER BY mr.checked_at DESC
			LIMIT 1
		) last_result ON TRUE
		WHERE m.is_active = TRUE
		AND (
			last_result.checked_at IS NULL
			OR last_result.checked_at <=
				NOW() - (
					m.interval_seconds *
					INTERVAL '1 second'
				)
		)
		ORDER BY m.id ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	monitors := make(
		[]model.Monitor,
		0,
	)

	for rows.Next() {
		var monitor model.Monitor

		err := rows.Scan(
			&monitor.ID,
			&monitor.UserID,
			&monitor.Name,
			&monitor.URL,
			&monitor.MonitorType,
			&monitor.IntervalSeconds,
			&monitor.IsActive,
			&monitor.CreatedAt,
			&monitor.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		monitors = append(
			monitors,
			monitor,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return monitors, nil
}
