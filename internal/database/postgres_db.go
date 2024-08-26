package database

import (
	"context"
	"database/sql"
	"fmt"
	"poster/graph/model"
	"strconv"
	"strings"
)

type postgresDB struct {
	db *sql.DB
}

func (db *postgresDB) CloseDB() {
	db.db.Close()
}

func (db *postgresDB) CreatePost(ctx context.Context, newPost model.NewPostInput) (int, error) {
	var row sql.NullInt32
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	insertPostStmt := `INSERT INTO public.posts(title, author_id, body, is_commented) VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRow(insertPostStmt, newPost.Title, newPost.AuthorID, newPost.Body, newPost.IsCommented).Scan(&row)
	if err != nil {
		return 0, err
	}

	postId := row.Int32

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return int(postId), nil

}

func (db *postgresDB) UpdatePost(ctx context.Context, post model.UpdPostInput) (int, error) {
	var row sql.NullInt32
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var updFields strings.Builder

	if post.Title != nil {
		updFields.WriteString(", title = '" + *post.Title + "'")
	}

	if post.Body != nil {
		/*if len(updFields.String()) != 0 {
			updFields.WriteString(", ")
		}*/
		updFields.WriteString(", body = '" + *post.Body + "'")
	}

	if post.IsCommented != nil {
		/*if len(updFields.String()) != 0 {
			updFields.WriteString(", ")
		}*/
		updFields.WriteString(", is_commented = " + fmt.Sprint(*post.IsCommented))
	}

	updatePostStmt := `UPDATE public.posts
	                   SET upd_date = now()` + updFields.String() +
		` WHERE id = $1 
		       and is_deleted = false
		RETURNING id`

	err = tx.QueryRow(updatePostStmt, post.PostID).Scan(&row)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	postId := row.Int32

	return int(postId), nil
}

func (db *postgresDB) DeletePost(ctx context.Context, id int) (int, error) {
	var row sql.NullInt32
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	deletePostStmt := `UPDATE public.posts
	                   SET is_deleted = true, upd_date = now()
		WHERE id = $1 
		       and is_deleted = false
		RETURNING id`
	err = tx.QueryRow(deletePostStmt, id).Scan(&row)
	if err != nil {
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	postId := row.Int32

	return int(postId), nil
}

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

func (db *postgresDB) GetPosts(ctx context.Context, userId int) ([]*model.Post, error) {

	type postPG struct {
		id          sql.NullInt32
		title       sql.NullString
		authorId    sql.NullInt32
		body        sql.NullString
		createDate  sql.NullTime
		updDate     sql.NullTime
		isCommented sql.NullBool
		likesCnt    sql.NullInt32
		dislikesCnt sql.NullInt32
	}

	likesStmt := `SELECT p.id,
    				     p.title,
    					 p.author_id ,
    					 p.body,
    					 p.create_date,
    					 p.upd_date,
    					 p.is_commented,
						 sum(case when c.is_like = true then 1 else 0 end ) likes,
						 sum(case when c.is_like = false then 1 else 0 end ) dislikes
				  FROM public.posts p
				  LEFT JOIN public.post_likes c ON p.id = c.post_id
				  WHERE p.author_id  = $1
				        AND p.is_deleted = false
				  ORDER BY p.upd_date DESC
			      GROUP BY p.id,
    				     p.title,
    					 p.author_id ,
    					 p.body,
    					 p.create_date,
    					 p.upd_date,
    					 p.is_commented`
	rows, err := db.db.Query(likesStmt, userId)

	if err != nil {
		return nil, err
	}
	var postList []*model.Post
	for rows.Next() {
		var p postPG

		if err = rows.Scan(&p.id, &p.title, &p.authorId, &p.body, &p.createDate, &p.updDate, &p.isCommented, &p.likesCnt, &p.dislikesCnt); err != nil {
			return nil, err
		}

		likes := int(p.likesCnt.Int32)
		dislikes := int(p.dislikesCnt.Int32)
		postList = append(postList, &model.Post{
			ID:          strconv.Itoa(int(p.id.Int32)),
			Title:       p.title.String,
			AuthorID:    strconv.Itoa(int(p.authorId.Int32)),
			Body:        p.body.String,
			CreateDate:  p.createDate.Time,
			UpdDate:     p.updDate.Time,
			IsCommented: p.isCommented.Bool,
			Likes:       &likes,
			Dislikes:    &dislikes,
		})
	}

	return postList, nil
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
