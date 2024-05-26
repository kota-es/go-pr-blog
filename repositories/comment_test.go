package repositories_test

import (
	"myapi/models"
	"myapi/repositories"
	"testing"
)

func TestSelectCommentList(t *testing.T) {
	articleID := 1
	got, err := repositories.SelectCommentList(testDB, articleID)
	if err != nil {
		t.Fatal(err)
	}

	for _, comment := range got {
		if comment.ArticleID != articleID {
			t.Errorf("want comment of articleID %d but got ID %d\n", articleID, comment.ArticleID)
		}
	}
}

func TestInsertComment(t *testing.T) {
	comment := models.Comment{
		ArticleID: 1,
		Message:   "CommentInsertTest",
	}

	before, err := repositories.SelectCommentList(testDB, comment.ArticleID)
	if err != nil {
		t.Fatal(err)
	}

	newComment, err := repositories.InsertComment(testDB, comment)
	if err != nil {
		t.Error(err)
	}

	after, err := repositories.SelectCommentList(testDB, comment.ArticleID)
	if err != nil {
		t.Fatal(err)
	}

	if len(after) != len(before)+1 {
		t.Errorf("Comment is not inserted\n")
	}

	t.Cleanup(func() {
		const sqlStr = "DELETE FROM comments WHERE comment_id = ?;"
		testDB.Exec(sqlStr, newComment.CommentID)
	})
}
