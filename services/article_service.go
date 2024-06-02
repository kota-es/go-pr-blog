package services

import (
	"database/sql"
	"errors"
	"myapi/apperrors"
	"myapi/models"
	"myapi/repositories"
)

func (s *MyAppService) PostArticleService(article models.Article) (models.Article, error) {
	newArticle, err := repositories.InsertArticle(s.db, article)
	if err != nil {
		err = apperrors.InsertDataFailed.Wrap(err, "fail to record data")
		return models.Article{}, err
	}
	return newArticle, nil
}

func (s *MyAppService) GetArticleListService(page int) ([]models.Article, error) {
	articleList, err := repositories.SelectArticleList(s.db, page)
	if err != nil {
		err = apperrors.GetDataFailed.Wrap(err, "fail to get data")
		return nil, err
	}

	if len(articleList) == 0 {
		err := apperrors.NAData.Wrap(ErrNoData, "no data")
		return nil, err
	}

	return articleList, nil
}

func (s *MyAppService) GetArticleService(articleID int) (models.Article, error) {
	var article models.Article
	var commentList []models.Comment
	var articleGetErr, commentGetErr error

	type artcleResult struct {
		article models.Article
		err     error
	}
	articleData := make(chan artcleResult)
	defer close(articleData)

	go func(ch chan<- artcleResult, db *sql.DB, articleID int) {
		article, err := repositories.SelectArticleDetail(db, articleID)
		ch <- artcleResult{article: article, err: err}
	}(articleData, s.db, articleID)

	type commentResult struct {
		commentList *[]models.Comment
		err         error
	}
	commentData := make(chan commentResult)
	defer close(commentData)

	go func(ch chan<- commentResult, db *sql.DB, articleID int) {
		commentList, err := repositories.SelectCommentList(db, articleID)
		ch <- commentResult{commentList: &commentList, err: err}
	}(commentData, s.db, articleID)

	for i := 0; i < 2; i++ {
		select {
		case ar := <-articleData:
			article, articleGetErr = ar.article, ar.err
		case cr := <-commentData:
			commentList, commentGetErr = *cr.commentList, cr.err
		}
	}

	if articleGetErr != nil {
		if errors.Is(articleGetErr, sql.ErrNoRows) {
			err := apperrors.NAData.Wrap(articleGetErr, "no data")
			return models.Article{}, err
		}
		err := apperrors.GetDataFailed.Wrap(articleGetErr, "fail to get data")
		return models.Article{}, err
	}

	if commentGetErr != nil {
		err := apperrors.GetDataFailed.Wrap(commentGetErr, "fail to get data")
		return models.Article{}, err
	}

	article.CommentList = append(article.CommentList, commentList...)

	return article, nil
}

func (s *MyAppService) PostNiceService(article models.Article) (models.Article, error) {
	err := repositories.UpdateNiceNum(s.db, article.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = apperrors.NoTargetData.Wrap(err, "does not exist target article")
			return models.Article{}, err
		}
		err = apperrors.UpdateDataFailed.Wrap(err, "fail to update nice count")
		return models.Article{}, err
	}

	return models.Article{
		ID:        article.ID,
		Title:     article.Title,
		Contents:  article.Contents,
		UserName:  article.UserName,
		NiceNum:   article.NiceNum + 1,
		CreatedAt: article.CreatedAt,
	}, nil
}
