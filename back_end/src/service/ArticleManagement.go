package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	//"rabbit_gather/src/handler"
	"rabbit_gather/src/neo4j_db"
	"rabbit_gather/util"
)

type ArticleManagement struct {
}

//func (a *ArticleManagement) GetHandler(handlerName auth.APIPermissionBitmask) gin.HandlerFunc {
//	switch handlerName {
//	case auth.PostArticle:
//		return a.PostArticleHandler
//	case auth.SearchArticle:
//		return a.SearchArticleHandler
//
//	default:
//		panic("No Such GetHandler")
//	}
//}

func (w *ArticleManagement) SearchArticleHandler(c *gin.Context) {
	type SearchArticleRequest struct {
		Position PositionStruct `json:"position"`
		Radius   int            `json:"radius"`
	}
	var searchArticleRequest SearchArticleRequest
	err := util.ParseRequestJson(c.Request.Body, &searchArticleRequest)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("SearchArticleHandler - parseRequestJson error : %s", err.Error())
		return
	}
}

type PositionStruct struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (w *ArticleManagement) PostArticleHandler(c *gin.Context) {
	articleReceived := struct {
		Title    string         `json:"title"`
		Content  string         `json:"content"`
		Position PositionStruct `json:"position"`
	}{}
	err := util.ParseRequestJson(c.Request.Body, &articleReceived)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		log.Printf("PostArticleHandler - parseRequestJson error : %s", err.Error())
		return
	}
	fmt.Println("Title : ", articleReceived.Title)
	fmt.Println("Content : ", articleReceived.Content)
	fmt.Println("Position : ", articleReceived.Position)
	res, err := neo4j_db.RunScriptWithScriptFile(
		"sql/create_new_article.cyp",
		map[string]interface{}{
			"username":  "A Name",
			"title":     articleReceived.Title,
			"content":   articleReceived.Content,
			"longitude": articleReceived.Position.Longitude,
			"latitude":  articleReceived.Position.Latitude,
		})
	if err != nil {
		panic("Error APIServer - PostArticleHandler : " + err.Error())
	}
	fmt.Println("neo4jTest - res :", res)
	c.JSON(200, gin.H{
		"result": articleReceived,
	})
}
