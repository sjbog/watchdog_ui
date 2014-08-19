package	security

import	(
	"encoding/json"
	"strings"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"
//	"watchdog_ui/app/controllers"
)

type	UserAuth	struct {

	HttpProtocol	string
	UserAgent		string
	RemoteAddr 		string
}

var	(
	AuthCache		= cache.NewInMemoryCache ( cache.DEFAULT )
	SessionCache	= cache.NewInMemoryCache ( cache.DEFAULT )
)

func	UserAuthGenerate	( request  * revel.Request )	( userAuth  * UserAuth )	{

	userAuth	= new ( UserAuth )

	userAuth.HttpProtocol	= request.Proto
	userAuth.UserAgent		= request.Header.Get ( "User-Agent" )

//	proxy1, proxy2, IP	OR just IP
	userAuth.RemoteAddr		= request.Header.Get ( "X-Forwarded-For" )

	if	userAuth.RemoteAddr == ""	{
		userAuth.RemoteAddr		= strings.Split ( request.RemoteAddr, ":" ) [ 0 ]
	}	else	{
		userAuth.RemoteAddr		+= ", " + strings.Split ( request.RemoteAddr, ":" ) [ 0 ]
	}


	return
}


func	( self  * UserAuth )	Equal	( userAuth  * UserAuth )		bool	{
	return	self.HttpProtocol == userAuth.HttpProtocol	&&
		self.UserAgent	== userAuth.UserAgent	&&
		self.RemoteAddr	== userAuth.RemoteAddr
}

func	( self  * UserAuth )	ToString	()		( result  * string )	{
	result		= new ( string )
	v, _	:= json.Marshal ( self )
	* result	= string ( v )
	return
}

func	( self  * UserAuth )	FromString	( s  * string )		( result  * UserAuth )	{
	result	= new ( UserAuth )
	json.Unmarshal ( []byte ( * s ),  result )
	return
}

func	( self  * UserAuth )	Save ( Session  * revel.Session )		{
	AuthCache.Set ( Session.Id (), * self, cache.DEFAULT )
}

type	SessionData	map [ string ] interface{}

func	GetSessionData	( Session  * revel.Session )		( sessionData  * SessionData )	{
	sessionData	= & SessionData {}

	if	err := SessionCache.Get ( Session.Id (), sessionData ); err != nil	{
		sessionData	= & SessionData {}
	}
	return
}

func	( self  * SessionData )	Save ( Session  * revel.Session )		{
	SessionCache.Set ( Session.Id (), * self, cache.DEFAULT )
}
