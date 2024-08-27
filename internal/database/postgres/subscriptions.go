package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/KseniiaTD/poster/graph/model"
)

func (db *postgresDB) CreateSubscription(ctx context.Context, subscr model.SubscrInput) (*string, error) {

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var id sql.NullInt32

	insertPostStmt := `INSERT INTO public.subscriptions(user_id, post_id) VALUES ($1, $2)
	                   ON CONFLICT(user_id, post_id) DO UPDATE 
	                   SET is_deleted = false,
					       upd_date = now()
					   RETURNING id`
	err = tx.QueryRow(insertPostStmt, subscr.UserID, subscr.PostID).Scan(&id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	str := strconv.Itoa(int(id.Int32))
	return &str, nil
}

func (db *postgresDB) DeleteSubscription(ctx context.Context, subscr model.SubscrInput) (*string, error) {

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var id sql.NullInt32

	updPostStmt := `UPDATE public.subscriptions
	                   SET is_deleted = true,
					       upd_date = now()
					   WHERE user_id = $1 AND post_id = $2
					   RETURNING id`
	err = tx.QueryRow(updPostStmt, subscr.UserID, subscr.PostID).Scan(&id)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	str := strconv.Itoa(int(id.Int32))
	return &str, nil
}

func (db *postgresDB) CheckSubscription(ctx context.Context, subscr model.Subscr) error {
	var cnt sql.NullInt32

	subscrStmt := `SELECT COUNT(*)
			 FROM public.subscriptions s
			 WHERE user_id = $1 AND post_id = $2 AND is_deleted = false`
	err := db.db.QueryRow(subscrStmt, subscr.UserID, subscr.PostID).Scan(&cnt)

	if err != nil {
		return err
	}

	if cnt.Int32 == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

func (db *postgresDB) GetCntNewCommentsForSubscriber(ctx context.Context, subscr model.Subscr) (int, error) {
	var cnt sql.NullInt32

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	deletePostStmt := `DELETE FROM public.new_comments
	                   WHERE subscription_id = (SELECT s.id FROM public.subscriptions s
	              WHERE s.user_id = $1 
			        AND s.post_id = $2 
			        AND s.is_deleted = false) 
					   RETURNING *`
	err = tx.QueryRow(deletePostStmt, subscr.UserID, subscr.PostID).Scan(&cnt)

	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return int(cnt.Int32), nil
}
