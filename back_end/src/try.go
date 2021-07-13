package main

import "rabbit_gather/src/logger"

type ABC string

func (a ABC) name() string {
	return string(a)
}

// API list
type PermissionCode uint32
type APIPermissionCode uint32

const (
	OK                              = APIPermissionCode(^uint32(0))
	Admin                           = APIPermissionCode(uint32(0))
	SearchArticle APIPermissionCode = 1 << iota // token is malformed
	PostArticle
	Login
	GetPeerID
	PeerWebsockt
	PeerWebsocketHandler
	GetPeerIDHandler
	UserAccountLoginHandler
	PermissionCheck
)

//const (
//	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
//	Ltime                         // the time in the local time zone: 01:23:23
//	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
//	Llongfile                     // full file name and line number: /a/b/c/d.go:23
//	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
//	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
//	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
//	LstdFlags     = Ldate | Ltime // initial values for the standard logger
//)

func main() {
	log := logger.NewLogger("")
	log.ERROR.Println("Error message.")
	log.WARNING.Println("WARNING message.")
	log.DEBUG.Println("DEBUG message.")
	//errorLogger := log.New(os.Stdout,"ERROR: ", log.Ltime|log.Ldate|log.Llongfile|log.Lmsgprefix)
	//errorLogger.Println("123ABC")
	//warningLogger := log.New(os.Stdout,"WARNING: ", log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix)
	//warningLogger.Println("123ABC")
	//debugLogger := log.New(os.Stdout,"DEBUG: ", log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix)
	//debugLogger.Println("123ABC")
	//
	//
	//loggerInst1 := log.New(os.Stdout,"---", log.Ltime|log.Lshortfile)
	//loggerInst1.Println("123ABC")
	//fmt.Println(OK)
	//fmt.Println(SearchArticle)
	////fmt.Println(OK)
	////fmt.Println(OK)
	////fmt.Println(OK)
	//userPermission := SearchArticle | PostArticle
	//fmt.Println(userPermission)
	//fmt.Println(SearchArticle)
	//fmt.Println(PostArticle)
	//fmt.Println(Login)
	//fmt.Println("----------")
	//fmt.Println(SearchArticle & userPermission)
	//fmt.Println(PostArticle & userPermission)
	//fmt.Println(OK & userPermission)
	//fmt.Println(Login & userPermission)     //0
	//fmt.Println(GetPeerID & userPermission) //0
	//fmt.Println(Admin & userPermission)     //0

	//if 1&1{
	//}
	//var Errors =  uint32(1)
	//fmt.Println(1 << 0)
	//fmt.Println(jwt.ValidationErrorMalformed )
	//fmt.Println(Errors)
	//fmt.Println(Errors&jwt.ValidationErrorMalformed )
	//fmt.Println(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet )
	//fmt.Println(1 &(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) )
	//
	////fmt.Println(jwt.ValidationErrorMalformed )
	////fmt.Println(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet)
	//if Errors&jwt.ValidationErrorMalformed != 0 {
	//	fmt.Println("That's not even a token")
	//} else if Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
	//	fmt.Println("Timing is everything")
	//} else {
	//	fmt.Println("Couldn't handle this token:")
	//}

	//a := ABC("456")
	//fmt.Println(a.name())
	//fmt.Println("4444444")
	//theuuid := uuid.New()
	//a := [16]byte(uuid.New())
	//p ,_:=filepath.Abs("ssl/crt/meowalien_com.crt")
	//fmt.Println("ssl/crt/meowalien_com.crt  --  ",p) // F:\GoTest\GoTest\master.exe <nil>
	//fmt.Println(string(a[:]))
	//fmt.Println(string(uuid.New().NodeID()))
}
