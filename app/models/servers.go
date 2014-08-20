package	models

import	(
	"fmt"
	"errors"
	"strconv"
	"strings"

	"github.com/revel/revel"
	robfig_config	"github.com/robfig/config"
)

const	SERVERS_CONF		= "servers.conf"
const	CMD_SECTION_SUFFIX	= "/commands"

type Servers	map [ string ] * Server


func LoadServers ()	( servers  * Servers, err  error )	{
	servers	= new ( Servers );	* servers = make ( Servers )
	MergedConfig, err	:= revel.LoadConfig ( SERVERS_CONF )

	if	err != nil	{
		if	err.Error ()	== "not found"	{
			err	= nil
			return
		}
		err	= errors.New ( fmt.Sprintf ( "Could not load config file \"%s\" : %s", SERVERS_CONF, err.Error () ) )
		return
	}

	for	_, sectionName	:= range ( MergedConfig .Raw () .Sections () )	{

		if	sectionName == "DEFAULT"	|| strings.HasSuffix ( sectionName, CMD_SECTION_SUFFIX )	{
			continue
		}
		MergedConfig .SetSection ( sectionName )

		server	:= NewServerFromConfig ( MergedConfig )
		( * server ).Label	= sectionName

		if	MergedConfig.HasSection ( sectionName + CMD_SECTION_SUFFIX )	{
			MergedConfig .SetSection ( sectionName + CMD_SECTION_SUFFIX )
			( * server ).Commands	= * LoadOptionsFromConfig ( MergedConfig )
		}

		( * servers ) [ sectionName ]	= server
	}
	return
}


func ( self  * Servers )	Save ()	( error )	{
	var config	= revel.NewEmptyConfig ()
	* config .Raw ()	= * robfig_config.New ( robfig_config.DEFAULT_COMMENT, robfig_config.ALTERNATIVE_SEPARATOR, true, true )

	for	_, server	:= range * self	{
		config.SetSection ( server.Label )
		config.SetOption ( "host", server.Host )
		config.SetOption ( "port", server.Port )
		config.SetOption ( "username", server.Username )
		config.SetOption ( "password", server.Password )
		config.SetOption ( "private_key", server.PrivateKeyPath )
		config.SetOption ( "query_interval", strconv.Itoa ( server.QueryIntervalSec ) )


		config.SetSection ( server.Label + "/commands" )

		for	cmd_name, cmd	:= range server.Commands	{
			config.SetOption ( cmd_name, cmd )
		}
	}
	return	config .Raw () .WriteFile ( revel.BasePath + "/conf/" + SERVERS_CONF, 0644, "remote servers" )
}


func ( self  * Servers )	Run ()	{
	for	_, x := range ( * self )	{
		x.Run ()
	}
}
func ( self  * Servers )	Start ()	{
	for	_, x := range ( * self )	{
		x.Start ()
	}
}
func ( self  * Servers )	Stop ()	{
	for	_, x := range ( * self )	{
		x.Stop ()
	}
}