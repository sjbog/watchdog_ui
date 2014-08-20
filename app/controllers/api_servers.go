package controllers

import	(
	"fmt"
	"strings"
	"io"
	"io/ioutil"
	"encoding/json"

	"watchdog_ui/app/models"
	"watchdog_ui/app/routes"

	"github.com/revel/revel"
//	"github.com/kr/prettsy"	//	fmt.Printf("%# v", pretty.Formatter( value ))
)

type ApiServers struct {
	*revel.Controller
}
//	TODO: don't show server's passwords for users

func (self *ApiServers) All () revel.Result {
	var	action	= strings.ToLower ( self.Params.Get ( "action" ) )
	switch	action	{
		case "save"	: ServersLastError	= ServersMap.Save ()

		case "reload"	:
			if	ServersLastError != nil	{
				return	self.RenderJson ( GenerateJsonStruct ( "", ServersLastError.Error () ) )
			}
			for	_, s	:= range ( * ServersMap )	{
				s.Delete ()
			}
			ServersMap, ServersLastError	= models .LoadServers ()
	}
	if	ServersLastError != nil	{
		revel.ERROR.Print ( ServersLastError )
		return	self.RenderJson ( GenerateJsonStruct ( "", ServersLastError.Error () ) )
	}

	PerformActions ( self.Params, ServersMap )
	return	self.RenderJson ( ServersMap )
}


func ( self  * ApiServers )	Show ( id  string )		( revel.Result )	{
	var server, result	= self.GetResource ( id )
	if	server == nil	{
		return	result
	}

	PerformActions ( self.Params, server )
	return	self.RenderJson ( server )
}

func ( self  * ApiServers )	GetResource ( id  string )		( * models.Server, revel.Result )	{

	if	ServersLastError != nil	{
		return	nil, self.RenderJson ( GenerateJsonStruct ( "", ServersLastError.Error () ) )
	}

	server, isFound	:= ( * ServersMap )[ id ]
	if	isFound == false	{
		self.Response.Status	= 404
		return	nil, self.RenderJson ( GenerateJsonStruct ( "", fmt.Sprintf ( "ID '%s' not found", id ) ) )
	}
	return	server, nil
}

func ( self  * ApiServers )	CreateResource ()	( * models.Server, revel.Result )	{
	var server	models.Server
	err	:= DecodeJsonPayload ( self.Request.Body, & server )

	if	err != nil	{
		return	nil,
			self.RenderJson ( GenerateJsonStruct ( "", fmt.Sprintf ( "Couldn't parse provided data : '%s'", err ) ) )
	}

	if	server.Label == ""	{
		return	nil, self.RenderJson ( GenerateJsonStruct ( "", "Cannot create unLabeled server" ) )
	}

	return	& server, nil
}

func ( self  * ApiServers )	Create ()		( revel.Result )	{
//	Create server instance from request.Body
	server, result	:= self.CreateResource ()
	if	server == nil	{
		return	result
	}
	( * ServersMap )[ server.Label ]	= server

	var returnData	= * GenerateJsonStruct ( fmt.Sprintf ( "New server \"%s\" was created", server.Label ), "" )
	returnData [ "url" ]	= routes.ApiServers.Show ( server.Label )

	return	self.RenderJson ( returnData )
}

func ( self  * ApiServers )	Alter ( id, method  string )	( revel.Result )	{
//	Check if server with given ID exists
	var server, result	= self.GetResource ( id )
	if	server == nil	{
		return	result
	}

	if	method == ""	{
		method	= self.Params.Get ( "_method" )
	}
	var returnData	= * GenerateJsonStruct ( "", "" )

	switch	strings.ToLower ( method )	{

		case "delete"	:
			returnData [ "result" ]	= fmt.Sprintf ( "%s was deleted", server.Label )
			delete ( * ServersMap, server.Label )
			server.Delete ()

//		PUT, POST, PATCH will replace / update
		default	:
			var	old_label	= server.Label
//			Create server instance from request.Body
			server, result	= self.CreateResource ()
			if	server == nil	{
				return	result
			}
			( * ServersMap )[ old_label ].Delete ()

			if	old_label != server.Label	{
				delete ( * ServersMap, old_label )
			}
			( * ServersMap )[ server.Label ]	= server

			returnData [ "result" ]	= fmt.Sprintf ( "%s was updated", old_label )
			returnData [ "url" ]	= routes.ApiServers.Show ( server.Label )
	}
	return	self.RenderJson ( returnData )
}

func PerformActions ( params  * revel.Params, serverInterface  models.ServerInterface )	{
	var action	= strings.ToLower ( params.Get ( "action" ) )

	switch	action	{
		case "start"	:
			serverInterface.Start ()
		case "stop"	:
			serverInterface.Stop ()
		case "run"	:
			serverInterface.Run ()
	}
}

func DecodeJsonPayload	( request_body  io.ReadCloser, v  interface {} )	( error )	{
    content, err	:= ioutil.ReadAll ( request_body )
    if err != nil {
        return err
    }

    err = json.Unmarshal(content, & v )
    if err != nil {
        return err
    }
    return nil
}

func GenerateJsonStruct	( response, error  string )	( * map [ string ] interface{} )	{
	return	& map [ string ] interface{} { "error" : error, "result" : response }
}