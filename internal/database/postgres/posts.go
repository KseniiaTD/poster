package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/KseniiaTD/poster/graph/model"
)

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
