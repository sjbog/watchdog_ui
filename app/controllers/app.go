package controllers

import	(
//	"fmt"
	"encoding/json"
	"watchdog_ui/app/models"
//	"watchdog_ui/app/jobs"
//	"github.com/revel/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"

//	"github.com/kr/pretty"
)

type App struct {
	*revel.Controller
}

var	(
	ServersMap		* models.Servers
	ServersLastError	error
)

func init ()	{
	revel.InterceptFunc ( CheckUserAuth, revel.BEFORE, revel.ALL_CONTROLLERS )

	revel.OnAppStart ( func ()	{
		ServersMap, ServersLastError	= models .LoadServers ()
		if	ServersLastError != nil	{
			revel.ERROR.Print ( ServersLastError )
		}
	})
}

func (self *App) Index() revel.Result {

	if	ServersLastError != nil	{
		self.RenderArgs [ "error" ]	= ServersLastError
	}	else	{
		var tmp,_	= json.Marshal ( ServersMap )
		self.RenderArgs [ "serversJSON" ]	= string ( tmp )
	}
	return self.Render ()
}