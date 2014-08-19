package controllers

import	(
	"fmt"
	"strings"
	"watchdog_ui/app/models"
	"watchdog_ui/app/routes"
	"github.com/revel/revel"

//	"github.com/kr/pretty"
)

type ApiServers struct {
	*revel.Controller
}

func (self *ApiServers) All () revel.Result {

	if	ServersLastError != nil	{
		var result	= map [ string ] interface{} { "error2" : ServersLastError, "error" : "a" }
		return	self.RenderJson ( result )
	}

	var	action	= strings.ToLower ( self.Params.Get ( "action" ) )
	switch	action	{
		case "save"	:
			ServersMap.Save ()
		case "reload"	:
			for	_, s	:= range ( * ServersMap )	{
				s.Delete ()
			}
			ServersMap, ServersLastError	= models .LoadServers ()
			if	ServersLastError != nil	{
				revel.ERROR.Print ( ServersLastError )
			}
	}

	PerformActions ( self.Params, ServersMap )
	return	self.RenderJson ( ServersMap )
}

func ( self  * ApiServers )	Show ( id  string )		( revel.Result )	{
	server, result	:= self.GetResource ( id )

	if	server == nil	{
		return	result
	}

	PerformActions ( self.Params, server )
	return	self.RenderJson ( server )
}

func ( self  * ApiServers )	GetResource ( id  string )		( * models.Server, revel.Result )	{

	if	ServersLastError != nil	{
		var result	= map [ string ] interface{} { "error" : ServersLastError }
		return	nil, self.RenderJson ( result )
	}

	server, isFound	:= ( * ServersMap )[ id ]

	if	isFound == false	{
		var result	= map [ string ] interface{} { "error" : fmt.Sprintf ( "ID '%s' not found", id ) }
		self.Response.Status	= 404
		return	nil, self.RenderJson ( result )
	}

	return	server, nil
}

func ( self  * ApiServers )	Create ()		( revel.Result )	{
	var	(
		returnData	= map [ string ] string {	"result" : "ok"	}
		server	= models.NewServerFromParams ( self.Params )
	)
	if	server.Label == ""	{
		returnData [ "error" ]	= "Cannot create unLabeled server"
		delete ( returnData, "result" )
		return	self.RenderJson ( returnData )
	}

	( * ServersMap )[ server.Label ]	= server
	returnData [ "url" ]	= routes.ApiServers.Show ( server.Label )

	return	self.RenderJson ( returnData )
}
func ( self  * ApiServers )	Action ( method, id  string )		( revel.Result )	{

	server, result	:= self.GetResource ( id )

	if	server == nil	{
		return	result
	}

	if	method == ""	{
		method	= self.Params.Get ( "_method" )
	}

	var returnData	= map [ string ] string {
		"result" : "ok",
	}

	switch	strings.ToLower ( method )	{
		case "delete"	:
			returnData [ "result" ]	= fmt.Sprintf ( "%s was deleted", server.Label )
			delete ( * ServersMap, server.Label )
			server.Delete ()

		case "update"	:	fallthrough
		default	:
			var updatedServer	= models.NewServerFromParams ( self.Params )
			if	updatedServer.Label == ""	{
				returnData [ "error" ]	= "Cannot create unLabeled server"
				delete ( returnData, "result" )
				return	self.RenderJson ( returnData )
			}
			var	oldLabel	= server.Label
			server.Delete ()
			if	oldLabel != updatedServer.Label	{
				delete ( * ServersMap, oldLabel )
			}
			server	= updatedServer
			( * ServersMap )[ server.Label ]	= server

			returnData [ "result" ]	= fmt.Sprintf ( "%s was updated", oldLabel )
			returnData [ "url" ]	= routes.ApiServers.Show ( server.Label )
	}

	return	self.RenderJson ( returnData )
}

func PerformActions ( params  * revel.Params, serverInterface  models.ServerInterface )	{
	var	(
		action	= strings.ToLower ( params.Get ( "action" ) )
	)
	switch	action	{
		case "start"	:
			serverInterface.Start ()
		case "stop"	:
			serverInterface.Stop ()
		case "run"	:
			serverInterface.Run ()
	}
}