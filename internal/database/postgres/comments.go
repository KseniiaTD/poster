package postgres

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/KseniiaTD/poster/graph/model"
)

func (db *postgresDB) CreateComment(ctx context.Context, newComment model.NewCommentInput) (int, error) {
	var row sql.NullInt32
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	insertCommentStmt := `INSERT INTO public.comments(parent_id, author_id, post_id, body) VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRow(insertCommentStmt, newComment.ParentID, newComment.AuthorID, newComment.PostID, newComment.Body).Scan(&row)
	if err != nil {
		return 0, err
	}

	commentId := row.Int32

	insertNewCommentStmt := `INSERT INTO public.new_comments(post_id, comment_id) VALUES ($1, $2)`
	_, err = tx.Query(insertNewCommentStmt, newComment.PostID, commentId)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return int(commentId), nil
}

func (db *postgresDB) UpdateComment(ctx context.Context, comment model.UpdCommentInput) (int, error) {
	var row sql.NullInt32
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	updateCommentStmt := `UPDATE public.comments
	                   SET body = $1, upd_date = now()
		WHERE id = $2 
		       and is_deleted = false
		RETURNING id`
	err = tx.QueryRow(updateCommentStmt, comment.Body, comment.ID).Scan(&row)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	commentId := row.Int32

	return int(commentId), nil
}

func (db *postgresDB) DeleteComment(ctx context.Context, id int) (int, error) {
	var row sql.NullInt32
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	deleteCommentStmt := `UPDATE public.comments
	                   SET is_deleted = true, upd_date = now()
		WHERE id = $1 
		       and is_deleted = false
		RETURNING id`
	err = tx.QueryRow(deleteCommentStmt, id).Scan(&row)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	commentId := row.Int32

	return int(commentId), nil
}

func (db *postgresDB) GetComments(ctx context.Context, postId int, commentId int, page int, perPage int) ([]*model.Comment, error) {
	type commentPG struct {
		id          sql.NullInt32
		body        sql.NullString
		authorId    sql.NullInt32
		postId      sql.NullInt32
		parentId    sql.NullInt32
		createDate  sql.NullTime
		updDate     sql.NullTime
		login       sql.NullString
		likesCnt    sql.NullInt32
		dislikesCnt sql.NullInt32
		commentsCnt sql.NullInt32
	}

	if perPage == 0 {
		perPage = 10
	}

	likesStmt := `WITH cnt_comments AS (
	SELECT parent_id, count(*) cnt 
	FROM comments 
	WHERE post_id = $1
	GROUP BY parent_id
)
,	parents AS (
	SELECT c.id,
	       case when c.is_deleted = true then 'comment was deleted' else c.body end body,
		   c.author_id,
		   c.post_id,
		   c.parent_id,
	       c.id as root_num,
		   u.login,
		   SUM(case when pl.is_like = true then 1 else 0 end) likes,
		   SUM(case when pl.is_like = false then 1 else 0 end) dislikes,
		   ROW_NUMBER() OVER (ORDER BY c.id) rn,
		   CASE WHEN c_child.cnt > 0 THEN c_child.cnt - 1 ELSE coalesce(c_child.cnt,0) END cnt,
		   c.create_date,
		   c.upd_date
	FROM comments c
	INNER JOIN users u ON c.author_id = u.id
	LEFT JOIN comment_likes pl ON c.id = pl.comment_id
	LEFT JOIN cnt_comments c_child ON c.id = c_child.parent_id 
	WHERE ((c.parent_id IS NULL AND $2 = 0)
	  OR c.parent_id = $2)
	  AND NOT (c.is_deleted = true AND COALESCE(c_child.cnt, 0) = 0)
	GROUP BY c.id,
	       c.body,
		   c.author_id,
		   c.post_id,
		   c.parent_id,
		   u.login,
	      c_child.cnt,
		  c.create_date,
		   c.upd_date
	)
, parents_perpage AS (
	SELECT c.id,
	       c.body,
		   c.author_id,
		   c.post_id,
		   c.parent_id,
	       c.root_num,
		   c.login,
		   c.likes,
		   c.dislikes,
		   c.cnt,
	       1 lvl,
		   c.create_date,
		   c.upd_date
	FROM parents c 
	WHERE rn > ($3::integer * $4::integer)
	LIMIT $3
)
, children AS (
	SELECT c.id,
	       case when c.is_deleted = true then 'comment was deleted' else c.body end body,
		   c.author_id,
		   c.post_id,
		   c.parent_id,
		   c.parent_id as root_num,
		   u.login,
		   sum(case when pl.is_like = true then 1 else 0 end) likes,
		   sum(case when pl.is_like = false then 1 else 0 end) dislikes,
		   coalesce(c_child.cnt,0) cnt,
	       2 lvl,
		   c.create_date,
		   c.upd_date
	FROM comments c
	INNER JOIN parents_perpage p ON c.parent_id = p.id 
	INNER JOIN users u ON c.author_id = u.id
	LEFT JOIN comment_likes pl ON c.id = pl.comment_id
	LEFT JOIN cnt_comments c_child ON c.id = c_child.parent_id
	WHERE NOT (c.is_deleted = true AND COALESCE(c_child.cnt, 0) = 0)
	GROUP BY c.id,
	       c.body,
		   c.author_id,
		   c.post_id,
		   c.parent_id,
		   u.login,
		   c.create_date,
		   c.upd_date,
	       c_child.cnt
   )
   , not_del_children AS (
   SELECT c.id,
	       c.body,
		   c.author_id,
		   c.post_id,
		   c.parent_id,
	       c.root_num,
		   c.login,
		   c.likes,
		   c.dislikes,
		   c.cnt,
	       c.lvl,
		   c.create_date,
		   c.upd_date,
		   min(c.id) over (partition by c.parent_id) min_id
	FROM children c
   )
   , posts AS (
     SELECT c.id,
	       c.body,
		   c.author_id,
		   c.post_id,
		   c.parent_id,
	       c.root_num,
		   c.login,
		   c.likes,
		   c.dislikes,
		   c.cnt,
		   c.lvl,
		   c.create_date,
		   c.upd_date
	FROM parents_perpage c
	UNION ALL
	SELECT c.id,
	       c.body,
		   c.author_id,
		   c.post_id,
		   c.parent_id,
	       c.root_num,
		   c.login,
		   c.likes,
		   c.dislikes,
		   c.cnt,
	       c.lvl,
		   c.create_date,
		   c.upd_date
	FROM not_del_children c
	WHERE c.min_id = c.id
   )
   SELECT  c.id,
	       c.body,
		   c.author_id,
		   c.post_id,
		   c.parent_id,
		   c.login,
		   c.create_date,
		   c.upd_date,
		   c.likes,
		   c.dislikes,
		   c.cnt
   FROM posts c
   ORDER BY root_num, lvl, id`
	rows, err := db.db.Query(likesStmt, postId, commentId, perPage, page)

	if err != nil {
		return nil, err
	}
	var commentList []*model.Comment
	for rows.Next() {
		var p commentPG

		if err = rows.Scan(&p.id,
			&p.body,
			&p.authorId,
			&p.postId,
			&p.parentId,
			&p.login,
			&p.createDate,
			&p.updDate,
			&p.likesCnt,
			&p.dislikesCnt,
			&p.commentsCnt,
		); err != nil {
			return nil, err
		}

		likes := int(p.likesCnt.Int32)
		dislikes := int(p.dislikesCnt.Int32)
		cnt := int(p.commentsCnt.Int32)
		parentId := strconv.Itoa(int(p.parentId.Int32))
		commentList = append(commentList, &model.Comment{
			ID:         strconv.Itoa(int(p.id.Int32)),
			ParentID:   &parentId,
			PostID:     strconv.Itoa(int(p.postId.Int32)),
			AuthorID:   strconv.Itoa(int(p.authorId.Int32)),
			Body:       p.body.String,
			Login:      p.login.String,
			CreateDate: p.createDate.Time,
			UpdDate:    p.updDate.Time,
			Likes:      &likes,
			Dislikes:   &dislikes,
			Cnt:        &cnt,
		})
	}

	return commentList, nil
}
