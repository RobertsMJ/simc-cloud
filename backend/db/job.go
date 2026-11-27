package db

import (
	"context"
	"log/slog"

	"github.com/RobertsMJ/simc-cloud/backend/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewJobRepository(log *slog.Logger, pool *pgxpool.Pool) *JobRepository {
	return &JobRepository{
		log:  log,
		pool: pool,
	}
}

func (r *JobRepository) Create(ctx context.Context, job models.Job) error {
	_, err := r.pool.Exec(ctx, `insert into jobs (
		id,
		status,
		spec
	) values ($1, $2, $3)`, job.ID, job.Status, job.Spec)
	return err
}

func (r *JobRepository) FindByID(ctx context.Context, id string) (*models.Job, error) {
	rows, err := r.pool.Query(ctx, `select id, status, spec, created_at, completed_at, simc_version, baseline_dps, baseline_dps_std_err, error from jobs where id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	job, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Job])

	return &job, err
}

func (r *JobRepository) RecordResult(ctx context.Context, result models.JobResult) (bool, error) {
	pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		tag, err := tx.Exec(
			ctx,
			`UPDATE jobs SET status='done', completed_at=now(),
			simc_version=$2, baseline_dps=$3, baseline_dps_std_err=$4
			WHERE id=$1 AND status='pending'`,
			result.JobID,
			result.SimcVersion,
			result.BaselineDPS,
			result.BaselineDPSStdErr,
		)
		if tag.RowsAffected() == 0 {
			return nil
		} // already recorded

		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{"gearset_results"},
			[]string{
				"job_id",
				"gearset_id",
				"loadout",
				"dps",
				"dps_std_err",
			},
			pgx.CopyFromSlice(len(result.Gearsets), func(i int) ([]interface{}, error) {
				gearset := result.Gearsets[i]
				return []interface{}{
					gearset.JobID,
					gearset.ID,
					gearset.Loadout,
					gearset.DPS,
					gearset.DPSStdErr,
				}, nil
			}),
		)

		return err
	})
	return true, nil
}
