package article_management

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"net/http"
	"rabbit_gather/src/neo4j_db"
	"rabbit_gather/src/server"
	"rabbit_gather/src/websocket/events"
	"rabbit_gather/util"
	"time"
)

type ArticleReceived struct {
	Title    string       `json:"title" form:"title"  binding:"required"`
	Content  string       `json:"content" form:"content"  binding:"required"`
	Position util.Point2D `json:"position" form:"position"  binding:"required"`
}

// PostArticleHandler create a new article
func (w *ArticleManagement) PostArticleHandler(c *gin.Context) {
	var articleReceived ArticleReceived
	err := c.ShouldBindJSON(&articleReceived)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"err": "input wrong",
		})
		log.DEBUG.Println("input wrong: ", err.Error())
		return
	}

	if !util.Is_WGS84_2D(articleReceived.Position.X, articleReceived.Position.Y) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input out of X,Y range. correct range : -90 < X <90 , -180 < X <180",
		})
		log.DEBUG.Println("wrong input out of X,Y range. correct range : -90 < X <90 , -180 < X <180")
		return
	}

	userClaim, err := server.ContextAnalyzer(c).ParseUserClaim()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "user claim error",
		})
		log.ERROR.Println("error ParseUserClaim:", err.Error())
		return
	}

	theID, err := insertArticleToDB(&articleReceived, userClaim.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.ERROR.Println("error when insertArticleToDB:", err.Error())
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"id": theID,
	})
	ArticleChangeBroker.Publish(&events.NewArticleEvent{
		Event:     events.NewArticle,
		Timestamp: time.Now().Unix(),
		ArticleID: theID,
		Position:  articleReceived.Position,
	})
}

// insertArticleToDB will insert a new article to RDBMS and store a "user
// create article at place" relationship to graph database (neo4j)
func insertArticleToDB(article *ArticleReceived, userID uint32) (int64, error) {
	tx, err := dbOperator.Begin()
	if err != nil {
		return 0, err
	}
	defer func(stmt *sql.Tx) {
		e := tx.Rollback()
		if e != sql.ErrTxDone && e != nil {
			log.ERROR.Println("error when Rollback: ", e.Error())
		}
	}(tx)

	insertToArticle := tx.Stmt(dbOperator.Statement("insert into `article` (title, content)\n    value (?, ?)\n;"))
	insertToArticleDetails := tx.Stmt(dbOperator.Statement("insert into `article_details` (article, coords)\n    value (?, Point(?, ?))\n;"))

	res, err := insertToArticle.Exec(article.Title, article.Content)
	if err != nil {
		return 0, err
	}
	articleID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	_, err = insertToArticleDetails.Exec(articleID, article.Position.X, article.Position.Y)
	if err != nil {
		return 0, err
	}
	log.TempLog().Println("HERE")
	session := neo4j_db.GetDriver().NewSession(neo4j.SessionConfig{})
	defer func(session neo4j.Session) {
		e := session.Close()
		if e != nil {
			log.ERROR.Println("error when Close neo4j_db Session: ", e.Error())
		}
	}(session)

	_, err = session.Run(
		util.GetFileStoredPlainText("sql/create_new_article.cyp"),
		map[string]interface{}{
			"user_id":    userID,
			"article_id": articleID,
			"x":          article.Position.X,
			"y":          article.Position.Y,
		},
	)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return articleID, nil
}
