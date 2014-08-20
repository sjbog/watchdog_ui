package controllers

import	(
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/revel/revel"

	"watchdog_ui/app/models"
	"watchdog_ui/app/routes"
	"watchdog_ui/app/security"
)

type User	struct {
	* revel.Controller
}

const	(
	MAX_LOGIN_TRIES					= 5
	UNSUCCESSFUL_LOGIN_TIMEOUT_MIN	= 60
)

func	( self  * User )	Login ( username, password string, remember_flag bool )	revel.Result	{

	var	sessionData	= * security.GetSessionData ( & self.Session )
	defer	sessionData.Save ( & self.Session )

	if	_, ok := sessionData [ "username" ] ;	ok	{
		return	self.Redirect ( ( * App ).Index )
	}

	self.Session.SetNoExpiration ()
	var	(
		user	= models.GetUser ( username )
		err		error
	)

//	hash, err	:= bcrypt.GenerateFromPassword ( [] byte ( password ), bcrypt.DefaultCost )
//	revel.INFO.Print ( string ( hash ), err )

	if	user != nil	{

		err	= bcrypt.CompareHashAndPassword ( user.HashedPassword, [] byte ( password ) )

		if	err == nil	{

			if	remember_flag	{
				self.Session.SetDefaultExpiration ()
			}

			security.UserAuthGenerate ( self.Request ) .Save ( & self.Session )
			sessionData [ "username" ]	= username

			return	self.Redirect ( routes.App.Index () )
		}
	}


	if	username != ""	&& password != ""	{
		self.RenderArgs [ "error" ]			= "Username or password is incorrect"

//		TODO : N tries left
//		if	_, ok := self.Session [ "loginTry" ] ;	ok	{
//			self.RenderArgs [ "warning" ]		= "N tries left"
//		}
	}


	self.Response.Out.Header ().Set ( "Requires-Auth", "1" )

	self.RenderArgs [ "username" ]		= username
	self.RenderArgs [ "remember_flag" ]	= remember_flag

	return	self.RenderTemplate ( "App/Login.html" )
}

func	CheckUserAuth ( controller  * revel.Controller )	revel.Result	{

	if	controller.Action == "Static.Serve"	||
		controller.Action == "App.Login"	||
		controller.Action == "User.Login"	||
		controller.Action == "User.Logout"	{

		return	nil
	}

	var	(
		userAuth		= new ( security.UserAuth )
		username		string
		sessionData		= * security.GetSessionData ( & controller.Session )
	)

	security.AuthCache.Get ( controller.Session.Id (), userAuth )

	if	v, ok	:= sessionData [ "username" ]; ok	{
		username	= v.( string )
	}

	if	userAuth != nil	&& username != ""	&& userAuth.Equal ( security.UserAuthGenerate ( controller.Request ) )	{
		return	nil
	}

	controller.Flash.Error ( "Please log in first" )
	controller.Response.Out.Header ().Set ( "Requires-Auth", "1" )
//	controller.Response.Status	= 401

	return	controller.Redirect ( ( * User ).Login )
}

func	( self  * User )	Logout ()	revel.Result	{

	for	key := range self.Session	{
		delete ( self.Session, key )
	}

	security.AuthCache.Delete ( self.Session.Id () )
	security.SessionCache.Delete ( self.Session.Id () )

	return	self.Redirect ( routes.App.Index () )
}



func	getUser ( self   * revel.Controller )	* models.User	{

//	CheckAuth ( self )
//	self	= self.( * revel.Controller )

	if	username, ok := self.Session [ "UserId" ]; ok	{
		return	models.GetUser ( username )
	}

	return	nil
}