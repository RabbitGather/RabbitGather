package article_management

import (
	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"net/http"
	"rabbit_gather/src/neo4j_db"
	"rabbit_gather/util"
)

type SearchArticleRequest struct {
	Position  PositionStruct `json:"position" form:"position"  binding:"required"`
	MinRadius float64        `json:"min_radius" form:"min_radius"`
	MaxRadius float64        `json:"max_radius" form:"max_radius" binding:"required"`
}

func (w *ArticleManagement) SearchArticleHandler(c *gin.Context) {
	var searchArticleRequest SearchArticleRequest
	err := c.ShouldBindQuery(&searchArticleRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err": "wrong input",
		})
		log.DEBUG.Println("wrong input: ", err.Error())
		return
	}
	log.TempLog().Println(pretty.Sprint(searchArticleRequest))
	log.TempLog().Println("searchArticleRequest.Position.X: ", searchArticleRequest.Position.X)
	log.TempLog().Println("25.040056717110396: ", 25.040056717110396)
	session := neo4j_db.GetDriver().NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(
		util.GetFileStoredPlainText("sql/search_article_with_radius.cyp"),
		//"MATCH  (position:Position)-[]-(article:Article)\nWITH  article ,position ,distance(position.pt, point({longitude:$longitude, latitude:$latitude})) as distance\n  WHERE distance>$min_radius AND distance< $max_radius\nRETURN position , article , distance;",
		map[string]interface{}{
			"longitude":  searchArticleRequest.Position.X,
			"latitude":   searchArticleRequest.Position.Y,
			"min_radius": searchArticleRequest.MinRadius,
			"max_radius": searchArticleRequest.MaxRadius,
		},
	)

	log.TempLog().Println("result: ", result.Record())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "server error",
		})
		log.ERROR.Println("error when neo4j_db script run:", err.Error())
		return
	}
	type Article struct {
		ID        int64          `json:"id"`
		Title     string         `json:"title"`
		Content   string         `json:"content"`
		Timestamp int64          `json:"timestamp"`
		Position  PositionStruct `json:"position"`
		Distance  float64        `json:"distance"`
	}
	resultArticles := []Article{}
	for result.Next() {
		record := result.Record()

		article := Article{
			Position: PositionStruct{},
		}
		art, exist := record.Get("article")
		if exist {
			articleProps := art.(neo4j.Node).Props
			if id, exist := articleProps["id"]; exist {
				//log.TempLog().Println("id: ",id)
				article.ID = id.(int64)
			}
			if title, exist := articleProps["title"]; exist {
				article.Title = title.(string)
			}
			if content, exist := articleProps["content"]; exist {
				article.Content = content.(string)
			}
			if timestamp, exist := articleProps["timestamp"]; exist {
				article.Timestamp = timestamp.(int64)
			}
		}

		position, exist := record.Get("position")
		if exist {
			positionProps := position.(neo4j.Node).Props
			if point, exist := positionProps["pt"]; exist {
				//log.TempLog().Println("point: ",point)
				//log.TempLog().Println("point.(neo4j.Point2D).X: ",point.(neo4j.Point2D).X)
				//log.TempLog().Println("point.(neo4j.Point2D).Y: ",point.(neo4j.Point2D).Y)

				article.Position.Y = point.(neo4j.Point2D).Y
				article.Position.X = point.(neo4j.Point2D).X
			}
		}

		distance, exist := record.Get("distance")
		if exist {
			article.Distance = distance.(float64)
		}

		resultArticles = append(resultArticles, article)
	}
	c.JSON(200, gin.H{
		"articles": resultArticles,
	})
}
