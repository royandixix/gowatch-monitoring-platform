package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/royandi/gowatch/backend/internal/model"
)

type MonitorResultRepository struct {
	db *pgxpool.Pool
}

func NewMonitorResultRepository(
	db *pgxpool.Pool,
) *MonitorResultRepository {
	return &MonitorResultRepository{
		db: db,
	}
}

func (r *MonitorResultRepository) Create(
	ctx context.Context,
	monitorID int64,
	result model.MonitorResult,
) error {
	checkedAt := time.Now()

	if result.CheckedAt != "" {
		parsedTime, err := time.Parse(
			time.RFC3339,
			result.CheckedAt,
		)

		if err == nil {
			checkedAt = parsedTime
		}
	}

	var httpStatus any

	if result.HTTPStatus > 0 {
		httpStatus =
			result.HTTPStatus
	}

	var errorMessage any

	if result.Error != "" {
		errorMessage =
			result.Error
	}

	_, err := r.db.Exec(
		ctx,
		`
			INSERT INTO monitor_results (
				monitor_id,
				status,
				http_status,
				response_time_ms,
				error,
				checked_at
			)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
		monitorID,
		result.Status,
		httpStatus,
		result.ResponseTime,
		errorMessage,
		checkedAt,
	)

	return err
}

func (r *MonitorResultRepository) FindHistoryByUser(
	ctx context.Context,
	userID int64,
	filter model.MonitorHistoryFilter,
) ([]model.MonitorResultRecord, int64, error) {
	whereClause :=
		strings.Builder{}

	whereClause.WriteString(
		" WHERE m.user_id = $1 ",
	)

	args := []any{
		userID,
	}

	if filter.MonitorID > 0 {
		args = append(
			args,
			filter.MonitorID,
		)

		whereClause.WriteString(
			fmt.Sprintf(
				" AND m.id = $%d ",
				len(args),
			),
		)
	}

	if filter.Status != "" {
		args = append(
			args,
			filter.Status,
		)

		whereClause.WriteString(
			fmt.Sprintf(
				" AND mr.status = $%d ",
				len(args),
			),
		)
	}

	countQuery := `
		SELECT COUNT(*)
		FROM monitor_results mr
		INNER JOIN monitors m
			ON m.id = mr.monitor_id
	` + whereClause.String()

	var total int64

	err := r.db.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(
		&total,
	)

	if err != nil {
		return nil, 0, err
	}

	page :=
		filter.Page

	if page < 1 {
		page = 1
	}

	perPage :=
		filter.PerPage

	if perPage < 1 {
		perPage = 20
	}

	offset :=
		(page - 1) *
			perPage

	queryArgs :=
		append(
			[]any{},
			args...,
		)

	queryArgs = append(
		queryArgs,
		perPage,
	)

	limitPosition :=
		len(queryArgs)

	queryArgs = append(
		queryArgs,
		offset,
	)

	offsetPosition :=
		len(queryArgs)

	query := `
		SELECT
			mr.id,
			mr.monitor_id,
			m.name,
			m.url,
			mr.status,
			COALESCE(mr.http_status, 0),
			COALESCE(mr.response_time_ms, 0),
			COALESCE(mr.error, ''),
			mr.checked_at
		FROM monitor_results mr
		INNER JOIN monitors m
			ON m.id = mr.monitor_id
	` +
		whereClause.String() +
		fmt.Sprintf(
			`
				ORDER BY
					mr.checked_at DESC,
					mr.id DESC
				LIMIT $%d
				OFFSET $%d
			`,
			limitPosition,
			offsetPosition,
		)

	rows, err :=
		r.db.Query(
			ctx,
			query,
			queryArgs...,
		)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	results :=
		make(
			[]model.MonitorResultRecord,
			0,
		)

	for rows.Next() {
		var result model.MonitorResultRecord

		err := rows.Scan(
			&result.ID,
			&result.MonitorID,
			&result.MonitorName,
			&result.MonitorURL,
			&result.Status,
			&result.HTTPStatus,
			&result.ResponseTimeMS,
			&result.Error,
			&result.CheckedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		results = append(
			results,
			result,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *MonitorResultRepository) FindRecentByUser(
	ctx context.Context,
	userID int64,
	limit int,
) ([]model.MonitorResultRecord, error) {
	results, _, err :=
		r.FindHistoryByUser(
			ctx,
			userID,
			model.MonitorHistoryFilter{
				Page:    1,
				PerPage: limit,
			},
		)

	return results, err
}

func (r *MonitorResultRepository) FindByMonitorAndUser(
	ctx context.Context,
	monitorID int64,
	userID int64,
	limit int,
) ([]model.MonitorResultRecord, error) {
	results, _, err :=
		r.FindHistoryByUser(
			ctx,
			userID,
			model.MonitorHistoryFilter{
				MonitorID: monitorID,
				Page:      1,
				PerPage:   limit,
			},
		)

	return results, err
}

func (r *MonitorResultRepository) GetStatsByMonitorAndUser(
	ctx context.Context,
	monitorID int64,
	userID int64,
) (model.MonitorStats, error) {
	query := `
		WITH target_monitor AS (
			SELECT id
			FROM monitors
			WHERE id = $1
			AND user_id = $2
		),
		aggregate_stats AS (
			SELECT
				tm.id AS monitor_id,

				COUNT(mr.id)::BIGINT
					AS total_checks,

				COUNT(mr.id)
					FILTER (
						WHERE mr.status = 'UP'
					)::BIGINT
					AS total_up,

				COUNT(mr.id)
					FILTER (
						WHERE mr.status = 'DOWN'
					)::BIGINT
					AS total_down,

				ROUND(
					COALESCE(
						AVG(
							NULLIF(
								mr.response_time_ms,
								0
							)
						),
						0
					)
				)::BIGINT
					AS average_response_time_ms

			FROM target_monitor tm

			LEFT JOIN monitor_results mr
				ON mr.monitor_id = tm.id

			GROUP BY tm.id
		),
		latest_result AS (
			SELECT
				mr.monitor_id,
				mr.status,
				COALESCE(
					mr.http_status,
					0
				) AS http_status,
				COALESCE(
					mr.response_time_ms,
					0
				) AS response_time_ms,
				mr.checked_at

			FROM monitor_results mr

			INNER JOIN target_monitor tm
				ON tm.id = mr.monitor_id

			ORDER BY
				mr.checked_at DESC,
				mr.id DESC

			LIMIT 1
		)

		SELECT
			tm.id,

			stats.total_checks,

			stats.total_up,

			stats.total_down,

			CASE
				WHEN stats.total_checks = 0
					THEN 0::DOUBLE PRECISION
				ELSE
					ROUND(
						(
							stats.total_up::NUMERIC /
							stats.total_checks::NUMERIC
						) * 100,
						2
					)::DOUBLE PRECISION
			END AS uptime_percentage,

			stats.average_response_time_ms,

			COALESCE(
				latest.response_time_ms,
				0
			),

			COALESCE(
				latest.status,
				''
			),

			COALESCE(
				latest.http_status,
				0
			),

			latest.checked_at

		FROM target_monitor tm

		INNER JOIN aggregate_stats stats
			ON stats.monitor_id = tm.id

		LEFT JOIN latest_result latest
			ON latest.monitor_id = tm.id
	`

	var stats model.MonitorStats

	var lastCheckedAt pgtype.Timestamptz

	err := r.db.QueryRow(
		ctx,
		query,
		monitorID,
		userID,
	).Scan(
		&stats.MonitorID,
		&stats.TotalChecks,
		&stats.TotalUp,
		&stats.TotalDown,
		&stats.UptimePercentage,
		&stats.AverageResponseTimeMS,
		&stats.LastResponseTimeMS,
		&stats.LastStatus,
		&stats.LastHTTPStatus,
		&lastCheckedAt,
	)

	if err != nil {
		return model.MonitorStats{},
			err
	}

	if lastCheckedAt.Valid {
		value :=
			lastCheckedAt.Time

		stats.LastCheckedAt =
			&value
	}

	return stats, nil
}
