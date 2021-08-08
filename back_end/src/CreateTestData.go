package main

//
//
//
////var driver = neo4j_db.GetDriver()
//
//type Config struct {
//	DBUri    string `json:"DBUri"`
//	Username string `json:"Username"`
//	Password string `json:"Password"`
//}
//
//var config Config
//
//func init() {
//	err := util.ParseJsonConfic(&config, "config/neo4j_db.config.json")
//	if err != nil {
//		panic(err.Error())
//	}
//}
//func main() {
//	wc := sync.WaitGroup{}
//	var allCount int64 = 0
//	runner := func(ii int) {
//		driver, err := neo4j.NewDriver(config.DBUri, neo4j.BasicAuth(config.Username, config.Password, ""))
//		if err != nil {
//			panic(err.Error())
//		}
//		section := driver.NewSession(neo4j.SessionConfig{})
//		defer section.Close()
//		for i := 0; i < 1000000; i++ {
//			atomic.AddInt64(&allCount, 1)
//			fmt.Printf("%d - ID: %d\n", ii, allCount)
//			newArticle(section)
//		}
//		wc.Done()
//	}
//
//	//wc.Add(10)
//	for i := 0; i < 400; i++ {
//		wc.Add(i)
//		go runner(i)
//	}
//	wc.Wait()
//}
//
////var the_session  = driver.NewSession(neo4j.SessionConfig{})
//
//func newArticle(session neo4j.Session) {
//	username := util.RandStringRunes(rand.Intn(15))
//	theID := util.Snowflake().Int64()
//	title := util.RandStringRunes(rand.Intn(25))
//	content := util.RandStringRunes(rand.Intn(400))
//	positionX := randFloats(-180, 180)
//	positionY := randFloats(-90, 90)
//
//	_, err := session.Run(util.GetFileStoredPlainText("sql/create_new_article.cyp"),
//		map[string]interface{}{
//			"username":  username,
//			"title":     title,
//			"content":   content,
//			"longitude": positionX,
//			"latitude":  positionY,
//			"id":        theID,
//		},
//	)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}
