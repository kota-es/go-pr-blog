package repositories_test

import (
	"myapi/models"
	"myapi/repositories"
	"myapi/repositories/testdata"
	"testing"
)

func TestSelectArticleList(t *testing.T) {
	expectedNum := len(testdata.ArticleData)
	got, err := repositories.SelectArticleList(testDB, 1)
	if err != nil {
		t.Fatal(err)
	}

	if num := len(got); num != expectedNum {
		t.Errorf("got %d articles but want %d\n", num, expectedNum)
	}
}

func TestSelectArticleDetail(t *testing.T) {
	tests := []struct {
		testTitle string
		expected  models.Article
	}{
		{
			testTitle: "subtest1",
			expected:  testdata.ArticleData[0],
		},
		{
			testTitle: "subtest2",
			expected:  testdata.ArticleData[1],
		},
	}

	for _, test := range tests {
		t.Run(test.testTitle, func(t *testing.T) {
			got, err := repositories.SelectArticleDetail(testDB, test.expected.ID)
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != test.expected.ID {
				t.Errorf("ID: get %d but want %d\n", got.ID, test.expected.ID)
			}
			if got.Title != test.expected.Title {
				t.Errorf("Title: get %s but want %s\n", got.Title, test.expected.Title)
			}
			if got.Contents != test.expected.Contents {
				t.Errorf("Contents: get %s but want %s\n", got.Contents, test.expected.Contents)
			}
			if got.UserName != test.expected.UserName {
				t.Errorf("UserName: get %s but want %s\n", got.UserName, test.expected.UserName)
			}
			if got.NiceNum != test.expected.NiceNum {
				t.Errorf("NiceNum: get %d but want %d\n", got.NiceNum, test.expected.NiceNum)
			}
		})
	}
}

func TestInsertArticle(t *testing.T) {
	article := models.Article{
		Title:    "insertTest",
		Contents: "testest",
		UserName: "saki",
	}

	before, err := repositories.SelectArticleList(testDB, 1)
	if err != nil {
		t.Fatal("fail to get before data")
	}

	newArticle, err := repositories.InsertArticle(testDB, article)
	if err != nil {
		t.Error(err)
	}

	after, err := repositories.SelectArticleList(testDB, 1)
	if err != nil {
		t.Fatal("fail to get after data")
	}

	expectedArticleNum := len(before) + 1

	if len(after) != expectedArticleNum {
		t.Errorf("Article is not inserted\n")
	}

	t.Cleanup(func() {
		const sqlStr = "DELETE FROM articles WHERE article_id = ?;"

		testDB.Exec(sqlStr, newArticle.ID)
	})
}

func TestUpdateNiceNum(t *testing.T) {
	articleID := 1

	before, err := repositories.SelectArticleDetail(testDB, articleID)
	if err != nil {
		t.Fatal("fail to get before data")
	}

	err = repositories.UpdateNiceNum(testDB, articleID)
	if err != nil {
		t.Fatal(err)
	}

	after, err := repositories.SelectArticleDetail(testDB, articleID)
	if err != nil {
		t.Fatal("fail to get after data")
	}

	if before.NiceNum == after.NiceNum-1 {
		t.Errorf("NiceNum is not updated\n")
	}

	t.Cleanup(func() {
		const sqlStr = "UPDATE articles SET nice = ? WHERE article_id = ?;"

		testDB.Exec(sqlStr, before.NiceNum, articleID)
	})
}
