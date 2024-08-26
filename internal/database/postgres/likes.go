package postgres

import (
	"context"
	"database/sql"
)

func (db *postgresDB) UpdPostLikes(ctx context.Context, postId int, author_id int, isLike bool) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	likesStmt := `SELECT count(*), max(case when is_like= true then 1 else 0 end) isLike
			 FROM public.post_likes 
			 WHERE post_id = $1 and author_id = $2`
	rows, err := db.db.Query(likesStmt, postId, author_id)

	if err != nil {
		return err
	}

	var idSql sql.NullInt32
	var isLikeSql sql.NullBool
	for rows.Next() {
		if err = rows.Scan(&idSql, &isLikeSql); err != nil {
			return err
		}

	}

	if idSql.Int32 == 0 {
		insertStmt := `INSERT INTO public.post_likes(author_id, post_id, is_like) VALUES ($1, $2, $3)`
		_, err = tx.Exec(insertStmt, author_id, postId, isLike)

	} else if idSql.Int32 != 0 && isLikeSql.Bool == isLike {
		deleteStmt := `DELETE FROM public.post_likes
					   WHERE author_id = $1 AND post_id = $2 `
		_, err = tx.Exec(deleteStmt, author_id, postId)

	} else if idSql.Int32 != 0 && isLikeSql.Bool != isLike {
		updateStmt := `UPDATE public.post_likes
					   SET is_like = $3
					   WHERE author_id = $1 AND post_id = $2 `
		_, err = tx.Exec(updateStmt, author_id, postId, isLike)
	}

	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (db *postgresDB) UpdCommentLikes(ctx context.Context, author_id int, commentId int, isLike bool) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	likesStmt := `SELECT count(*), max(case when is_like= true then 1 else 0 end) isLike
			      FROM public.comment_likes 
			      WHERE comment_id = $1 AND author_id = $2`
	rows, err := db.db.Query(likesStmt, commentId, author_id)

	if err != nil {
		return err
	}

	var idSql sql.NullInt32
	var isLikeSql sql.NullBool
	for rows.Next() {
		if err = rows.Scan(&idSql, &isLikeSql); err != nil {
			return err
		}

	}

	if idSql.Int32 == 0 {
		insertStmt := `INSERT INTO public.comment_likes(author_id, comment_id, is_like) 
		               VALUES ($1, $2, $3)`
		_, err = tx.Exec(insertStmt, author_id, commentId, isLike)

	} else if idSql.Int32 != 0 && isLikeSql.Bool == isLike {
		deleteStmt := `DELETE FROM public.comment_likes
					   WHERE author_id = $1 AND comment_id = $2`
		_, err = tx.Exec(deleteStmt, author_id, commentId)

	} else if idSql.Int32 != 0 && isLikeSql.Bool != isLike {
		updateStmt := `UPDATE public.comment_likes
					   SET is_like = $3
					   WHERE author_id = $1 AND comment_id = $2`
		_, err = tx.Exec(updateStmt, author_id, commentId, isLike)
	}

	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (db *postgresDB) GetPostLikes(ctx context.Context, postId int) (int, int, error) {

	var likes, dislikes sql.NullInt32
	likesStmt := `SELECT SUM(case when p.is_like = true then 1 else 0 end) likes,
	                     SUM(case when p.is_like = false then 1 else 0 end) dislikes
			 FROM public.post_likes p
			 WHERE post_id = $1`
	rows, err := db.db.Query(likesStmt, postId)

	if err != nil {
		return 0, 0, err
	}

	for rows.Next() {
		if err = rows.Scan(&likes, &dislikes); err != nil {
			return 0, 0, err
		}
	}

	return int(likes.Int32), int(dislikes.Int32), nil
}

func (db *postgresDB) GetCommentLikes(ctx context.Context, commentId int) (int, int, error) {
	var likes, dislikes sql.NullInt32
	likesStmt := `SELECT SUM(case when p.is_like = true then 1 else 0 end) likes,
	                     SUM(case when p.is_like = false then 1 else 0 end) dislikes
			 FROM public.comment_likes p
			 WHERE comment_id = $1`
	rows, err := db.db.Query(likesStmt, commentId)

	if err != nil {
		return 0, 0, err
	}

	for rows.Next() {
		if err = rows.Scan(&likes, &dislikes); err != nil {
			return 0, 0, err
		}
	}

	return int(likes.Int32), int(dislikes.Int32), nil
}
