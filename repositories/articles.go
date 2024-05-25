package repositories

import (
	"database/sql"
	"myapi/models"

	_ "github.com/go-sql-driver/mysql"
)

const (
	articleNumPerPage = 5
)

func InsertArticle(db *sql.DB, article models.Article) (models.Article, error) {
	const sqlStr = "INSERT INTO articles (title, contents, username, nice, created_at) VALUES (?, ?, ?, 0, now());"

	var newArticle models.Article
	newArticle.Title, newArticle.Contents, newArticle.UserName = article.Title, article.Contents, article.UserName

	result, err := db.Exec(sqlStr, newArticle.Title, newArticle.Contents, newArticle.UserName)
	if err != nil {
		return models.Article{}, err
	}

	id, _ := result.LastInsertId()

	newArticle.ID = int(id)

	return newArticle, nil
}

func SelectArticleList(db *sql.DB, page int) ([]models.Article, error) {
	const sqlStr = "SELECT article_id, title, contents, username, nice FROM articles LIMIT ? OFFSET ?;"

	rows, err := db.Query(sqlStr, articleNumPerPage, (page-1)*articleNumPerPage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articleArray := make([]models.Article, 0)
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title, &article.Contents, &article.UserName, &article.NiceNum)
		if err != nil {
			return nil, err
		}
		articleArray = append(articleArray, article)
	}

	return articleArray, nil
}

func SelectArticleDetail(db *sql.DB, articleID int) (models.Article, error) {
	const sqlStr = "SELECT title, contents, username, nice, created_at FROM articles WHERE article_id = ?;"

	row := db.QueryRow(sqlStr, articleID)
	if row.Err() != nil {
		return models.Article{}, row.Err()
	}

	var article models.Article
	var createdTime sql.NullTime
	err := row.Scan(&article.Title, &article.Contents, &article.UserName, &article.NiceNum, &createdTime)
	if err != nil {
		return models.Article{}, err
	}

	if createdTime.Valid {
		article.CreatedAt = createdTime.Time
	}

	return article, nil
}

func UpdateNiceNum(db *sql.DB, articleID int) error {
	const sqlGetNice = "SELECT nice FROM articles WHERE article_id = ?"

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	row := tx.QueryRow(sqlGetNice, articleID)
	if row.Err() != nil {
		tx.Rollback()
		return row.Err()
	}

	var nicenum int
	err = row.Scan(&nicenum)
	if err != nil {
		tx.Rollback()
		return err
	}

	const sqlUpdateNice = "UPDATE articles SET nice = ? WHERE article_id = ?"

	_, err = tx.Exec(sqlUpdateNice, nicenum+1, articleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
