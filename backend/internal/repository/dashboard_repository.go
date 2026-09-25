package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/royandi/gowatch/backend/internal/model"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(
	db *pgxpool.Pool,
) *DashboardRepository {
	return &DashboardRepository{
		db: db,
	}
}

func percentage(
	success int64,
	total int64,
) float64 {
	if total <= 0 {
		return 0
	}

	value :=
		(float64(success) /
			float64(total)) *
			100

	return float64(
		int64(value*100+0.5),
	) / 100
}

func (r *DashboardRepository) GetDashboard(
	ctx context.Context,
	userID int64,
) (model.DashboardResponse, error) {
	var response model.DashboardResponse

	err :=
		r.loadMonitorSummary(
			ctx,
			userID,
			&response.Summary,
		)

	if err != nil {
		return model.DashboardResponse{},
			err
	}

	err =
		r.loadMonitoringSummary(
			ctx,
			userID,
			&response.Summary,
		)

	if err != nil {
		return model.DashboardResponse{},
			err
	}

	response.Monitors, err =
		r.loadMonitorStatuses(
			ctx,
			userID,
		)

	if err != nil {
		return model.DashboardResponse{},
			err
	}

	response.RecentActivity, err =
		r.loadRecentActivity(
			ctx,
			userID,
		)

	if err != nil {
		return model.DashboardResponse{},
			err
	}

	response.ResponseChart, err =
		r.loadResponseChart(
			ctx,
			userID,
		)

	if err != nil {
		return model.DashboardResponse{},
			err
	}

	return response, nil
}

func (r *DashboardRepository) loadMonitorSummary(
	ctx context.Context,
	userID int64,
	summary *model.DashboardSummary,
) error {
	query := `
		WITH latest_result AS (
			SELECT DISTINCT ON (
				mr.monitor_id
			)
				mr.monitor_id,
				mr.status

			FROM monitor_results mr

			INNER JOIN monitors monitor_owner
				ON monitor_owner.id = mr.monitor_id

			WHERE monitor_owner.user_id = $1

			ORDER BY
				mr.monitor_id,
				mr.checked_at DESC,
				mr.id DESC
		)

		SELECT
			COUNT(m.id)::BIGINT,

			COUNT(m.id)
				FILTER (
					WHERE
						m.is_active = TRUE
						AND latest.status = 'UP'
				)::BIGINT,

			COUNT(m.id)
				FILTER (
					WHERE
						m.is_active = TRUE
						AND latest.status = 'DOWN'
				)::BIGINT,

			COUNT(m.id)
				FILTER (
					WHERE
						m.is_active = FALSE
				)::BIGINT,

			COUNT(m.id)
				FILTER (
					WHERE
						m.is_active = TRUE
						AND (
							latest.status IS NULL
							OR latest.status NOT IN (
								'UP',
								'DOWN'
							)
						)
				)::BIGINT

		FROM monitors m

		LEFT JOIN latest_result latest
			ON latest.monitor_id = m.id

		WHERE m.user_id = $1
	`

	return r.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&summary.TotalMonitors,
		&summary.UpMonitors,
		&summary.DownMonitors,
		&summary.PausedMonitors,
		&summary.WaitingMonitors,
	)
}

func (r *DashboardRepository) loadMonitoringSummary(
	ctx context.Context,
	userID int64,
	summary *model.DashboardSummary,
) error {
	query := `
		SELECT
			COUNT(mr.id)
				FILTER (
					WHERE mr.checked_at >=
						NOW() - INTERVAL '24 hours'
				)::BIGINT
				AS total_checks_24h,

			COUNT(mr.id)
				FILTER (
					WHERE
						mr.checked_at >=
							NOW() - INTERVAL '24 hours'
						AND mr.status = 'UP'
				)::BIGINT
				AS total_up_24h,

			COUNT(mr.id)
				FILTER (
					WHERE mr.checked_at >=
						NOW() - INTERVAL '7 days'
				)::BIGINT
				AS total_checks_7d,

			COUNT(mr.id)
				FILTER (
					WHERE
						mr.checked_at >=
							NOW() - INTERVAL '7 days'
						AND mr.status = 'UP'
				)::BIGINT
				AS total_up_7d,

			ROUND(
				COALESCE(
					AVG(
						NULLIF(
							mr.response_time_ms,
							0
						)
					)
					FILTER (
						WHERE mr.checked_at >=
							NOW() - INTERVAL '24 hours'
					),
					0
				)
			)::BIGINT
				AS average_response_24h

		FROM monitor_results mr

		INNER JOIN monitors m
			ON m.id = mr.monitor_id

		WHERE m.user_id = $1
	`

	var up24H int64
	var up7D int64

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&summary.TotalChecks24H,
		&up24H,
		&summary.TotalChecks7D,
		&up7D,
		&summary.AverageResponseTime24H,
	)

	if err != nil {
		return err
	}

	summary.Uptime24HPercentage =
		percentage(
			up24H,
			summary.TotalChecks24H,
		)

	summary.Uptime7DPercentage =
		percentage(
			up7D,
			summary.TotalChecks7D,
		)

	return nil
}

func (r *DashboardRepository) loadMonitorStatuses(
	ctx context.Context,
	userID int64,
) ([]model.DashboardMonitorStatus, error) {
	query := `
		SELECT
			m.id,
			m.name,
			m.url,
			m.monitor_type,
			m.is_active,

			CASE
				WHEN m.is_active = FALSE
					THEN 'PAUSED'

				WHEN latest.status IS NULL
					THEN 'WAITING'

				ELSE latest.status
			END AS current_status,

			COALESCE(
				latest.http_status,
				0
			),

			COALESCE(
				latest.response_time_ms,
				0
			),

			latest.checked_at

		FROM monitors m

		LEFT JOIN LATERAL (
			SELECT
				mr.status,
				mr.http_status,
				mr.response_time_ms,
				mr.checked_at

			FROM monitor_results mr

			WHERE mr.monitor_id = m.id

			ORDER BY
				mr.checked_at DESC,
				mr.id DESC

			LIMIT 1
		) latest ON TRUE

		WHERE m.user_id = $1

		ORDER BY
			m.updated_at DESC,
			m.id DESC

		LIMIT 8
	`

	rows, err :=
		r.db.Query(
			ctx,
			query,
			userID,
		)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results :=
		make(
			[]model.DashboardMonitorStatus,
			0,
		)

	for rows.Next() {
		var result model.DashboardMonitorStatus

		var checkedAt pgtype.Timestamptz

		err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.URL,
			&result.MonitorType,
			&result.IsActive,
			&result.CurrentStatus,
			&result.LastHTTPStatus,
			&result.LastResponseTimeMS,
			&checkedAt,
		)

		if err != nil {
			return nil, err
		}

		if checkedAt.Valid {
			value :=
				checkedAt.Time

			result.LastCheckedAt =
				&value
		}

		results = append(
			results,
			result,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *DashboardRepository) loadRecentActivity(
	ctx context.Context,
	userID int64,
) ([]model.DashboardActivity, error) {
	query := `
		SELECT
			mr.id,
			mr.monitor_id,
			m.name,
			mr.status,
			COALESCE(
				mr.http_status,
				0
			),
			COALESCE(
				mr.response_time_ms,
				0
			),
			COALESCE(
				mr.error,
				''
			),
			mr.checked_at

		FROM monitor_results mr

		INNER JOIN monitors m
			ON m.id = mr.monitor_id

		WHERE m.user_id = $1

		ORDER BY
			mr.checked_at DESC,
			mr.id DESC

		LIMIT 10
	`

	rows, err :=
		r.db.Query(
			ctx,
			query,
			userID,
		)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results :=
		make(
			[]model.DashboardActivity,
			0,
		)

	for rows.Next() {
		var result model.DashboardActivity

		err := rows.Scan(
			&result.ID,
			&result.MonitorID,
			&result.MonitorName,
			&result.Status,
			&result.HTTPStatus,
			&result.ResponseTimeMS,
			&result.Error,
			&result.CheckedAt,
		)

		if err != nil {
			return nil, err
		}

		results = append(
			results,
			result,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *DashboardRepository) loadResponseChart(
	ctx context.Context,
	userID int64,
) ([]model.DashboardChartPoint, error) {
	query := `
		SELECT
			mr.id,
			mr.monitor_id,
			m.name,
			mr.status,
			COALESCE(
				mr.response_time_ms,
				0
			),
			mr.checked_at

		FROM monitor_results mr

		INNER JOIN monitors m
			ON m.id = mr.monitor_id

		WHERE
			m.user_id = $1
			AND mr.checked_at >=
				NOW() - INTERVAL '24 hours'

		ORDER BY
			mr.checked_at DESC,
			mr.id DESC

		LIMIT 48
	`

	rows, err :=
		r.db.Query(
			ctx,
			query,
			userID,
		)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	descending :=
		make(
			[]model.DashboardChartPoint,
			0,
		)

	for rows.Next() {
		var point model.DashboardChartPoint

		err := rows.Scan(
			&point.ID,
			&point.MonitorID,
			&point.MonitorName,
			&point.Status,
			&point.ResponseTimeMS,
			&point.CheckedAt,
		)

		if err != nil {
			return nil, err
		}

		descending = append(
			descending,
			point,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	results :=
		make(
			[]model.DashboardChartPoint,
			0,
			len(descending),
		)

	for index :=
		len(descending) - 1; index >= 0; index-- {
		results = append(
			results,
			descending[index],
		)
	}

	return results, nil
}
