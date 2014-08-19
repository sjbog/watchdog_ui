package	models

import	(
	"fmt"
	"errors"
	"strings"
	"bytes"
	"time"
	"strconv"
	"io/ioutil"

	"github.com/revel/revel"
	"code.google.com/p/go.crypto/ssh"
	"github.com/robfig/cron"
//	"github.com/kr/pretty"
)


const	SERVERS_CONF		= "servers.conf"
const	CMD_SECTION_SUFFIX	= "/commands"

type Server struct {
	Label			string
	PrivateKeyPath	string
//	PrivateKeyBytes	[] byte
	Username	string
	Password	string
	Host		string
	Port		string
	AuthMethods	[] ssh.AuthMethod	`json:"-"`
	ClientConnection	* ssh.Client	`json:"-"`
	QueryIntervalSec	int

	Schedule	cron.Schedule		`json:"-"`
	Cron		* cron.Cron			`json:"-"`

	Status	string

	Commands	map [ string ] string
	Responses	map [ string ] string
	Error		error
	ErrorMsg	string
}

type Servers	map [ string ] * Server

type ServerInterface	interface	{
	Run ()
	Start ()
	Stop ()
}

func LoadServers ()	( servers  * Servers, err  error )	{

//	servers	= nil
	MergedConfig, err	:= revel.LoadConfig ( SERVERS_CONF )

	if	err != nil	{
		err	= errors.New ( fmt.Sprintf ( "Could not load config file \"%s\" : %s", SERVERS_CONF, err.Error () ) )
		return
	}


	servers	= new ( Servers );	* servers = make ( Servers )

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

func NewServerFromConfig ( MergedConfig  * revel.MergedConfig )	( * Server )	{

	var	server	= Server {
		Username	: MergedConfig .StringDefault ( "username", "" ),
		Host	: MergedConfig .StringDefault ( "host", "" ),
		Port	: MergedConfig .StringDefault ( "port", "" ),
	}
	server.SetQueryInterval ( MergedConfig .IntDefault ( "query_interval", 60 ) )
	server.ParsePrivateKey ( MergedConfig .StringDefault ( "private_key", "" ) )

	if	password, isKeyPresent	:= MergedConfig .String ( "password" ); isKeyPresent	{
		server.SetPassword ( password )
	}
	return	& server
}

func NewServerFromParams ( Params  * revel.Params )	( server  * Server )	{
	server	= & Server {
		Label		: Params.Get ( "label" ),
		Username	: Params.Get ( "username" ),
		Host	: Params.Get ( "host" ),
		Port	: Params.Get ( "port" ),
	}
	queryInterval, err := strconv.Atoi ( Params.Get ( "query_interval" ) )

	if	err != nil	{
		server.SetQueryInterval ( 60 )
	}	else	{
		server.SetQueryInterval ( queryInterval )
	}
	server.ParsePrivateKey ( Params.Get ( "private_key" ) )
	server.SetPassword ( Params.Get ( "password" ) )

	var cmds	[][] string
	Params.Bind ( & cmds, "commands" )
	server.Commands	= make ( map [ string ] string )

	for	_, cmd	:= range ( cmds )	{
		if	len ( cmd ) == 2	&& cmd [ 0 ] != ""	&& cmd [ 1 ] != ""	{
			server.Commands [ cmd [ 0 ] ]	= cmd [ 1 ]
		}
	}
	return
}

func LoadOptionsFromConfig ( MergedConfig  * revel.MergedConfig )	( * map [ string ] string )	{

	var	options	= make ( map [ string ] string )

	for	_, optionName	:= range ( MergedConfig .Options ( "" ) )	{
		optionValue, _	:= MergedConfig .String ( optionName )
		if	optionName == ""	|| optionValue == ""	{
			continue
		}
		options [ optionName ]	= optionValue
	}
	return	& options
}

func ( self  * Server )	SetQueryInterval ( seconds  int )	{
	if	seconds < 3	{
		seconds	= 3
	}
	if	self.Cron != nil	{
		go	self.Cron.Stop ()
	}

	self.QueryIntervalSec	= seconds
	self.Schedule	= cron.Every ( time.Duration ( self.QueryIntervalSec ) * time.Second )
}

func ( self  * Server )	SetPassword ( password  string )	{
	if	password != ""	{
		self.Password	= password
		self.AuthMethods	= append ( self.AuthMethods, ssh.Password ( password ) )
	}
}

func ( self  * Server )	ParsePrivateKey ( filePath  string )	( error )	{
	if	filePath == ""	{
		return	nil
	}

	keyBytes, err	:= ioutil.ReadFile ( filePath )

	if	err != nil	{
		err			= errors.New ( fmt.Sprintf ( "Could not read ssh key \"%s\" : %s ", filePath, err.Error () ) )
		self.Error	= err
		self.ErrorMsg	= err.Error ()
		return	err
	}

	signer, err	:= ssh.ParsePrivateKey ( keyBytes )

	if	err != nil	{
		err			= errors.New ( fmt.Sprintf ( "Could not parse ssh key \"%s\" : %s ", filePath, err.Error () ) )
		self.Error	= err
		self.ErrorMsg	= err.Error ()
		return	err
	}

	self.PrivateKeyPath		= filePath
    self.AuthMethods	= append ( self.AuthMethods, ssh.PublicKeys ( signer ) )

	return	nil
}

func ( self  * Server )	Connect ()	( error )	{
	var	config	= & ssh.ClientConfig {
		User: self.Username,
		Auth: self.AuthMethods,
	}

	var	hostPort	= self.Host
	if	self.Port != ""	{
		hostPort	+= ":" + self.Port
	}

	ClientConnection, err	:= ssh.Dial ( "tcp", hostPort, config )

	if	err != nil {
		err	= errors.New ( "Failed to connect : " + err.Error () )
		self.ErrorMsg	= err.Error ()
	}	else	{
		self.ErrorMsg	= ""
	}
	self.ClientConnection	= ClientConnection
	self.Error	= err
	return	err
}

func ( self  * Server )	Query ( command  string )	( string )	{
//	Each ClientConn can support multiple interactive sessions, represented by a Session.
    session, err	:= self.ClientConnection.NewSession ()
    if	err != nil	{
        return	"Failed to create session: " + err.Error()
    }
    defer	session.Close()

    // Once a Session is created, you can execute a single command on the remote side using the Run method.
    var	StdOut, StdErr	bytes.Buffer
    session.Stdout	= & StdOut
    session.Stderr	= & StdErr

	var	response	= ""
    if	err := session.Run ( command ); err != nil {
        response	= "Failed to run: " + err.Error () + "; StdOUT : "
    }
	response	+= StdOut.String ()
	if	StdErr.Len () > 0	{
		response	+= "; StdERR : " + StdErr.String ()
	}
	return	response
}

func ( self  * Server )	Run ()	{
	if	self.ClientConnection == nil	&& self.Connect () != nil	{
		self.Stop ()
		return
	}
	if	self.Responses == nil	{
		self.Responses	= make ( map [string] string )
	}
//	TODO: improve with http://blog.golang.org/context
	for	label, cmd	:= range ( self.Commands )	{
		self.Responses [ label ]	= self.Query ( cmd )
	}
}

func ( self  * Server )	Start ()	{
	if	self.Cron != nil	{
		go	self.Cron.Stop ()
	}
	self.Cron	= cron.New ()
	self.Cron.Schedule ( self.Schedule, self )
	self.Cron.Start ()
	revel.INFO.Print ( "Starting " + self.Label )
	self.Status	= "running"
	self.Run ()
}

func ( self  * Server )	Continue ()	{
	if	self.Cron != nil	{
		self.Cron.Start ()
	}
}

func ( self  * Server )	Stop ()	{
	if	self.Cron != nil	{
		go self.Cron.Stop ()
	}
	if	self.ClientConnection != nil	{
		go	self.ClientConnection.Close ()
		self.ClientConnection	= nil
	}
	revel.INFO.Print ( "Stopping " + self.Label )
	self.Status	= "stopped"
}

func ( self  * Server )	Delete ()	{
	revel.INFO.Print ( "Deleting " + self.Label )
	self.Stop ()
	self.Cron	= nil
	self		= nil
}

//	No difference from delete and reassigned, since all fields are replaced

//func ( self  * Server )	Update ( other  * Server )	{
//	revel.INFO.Print ( "Updating " + self.Label )
//	self.Stop ()
//
//	if	other.Label != ""	&& self.Label	!= other.Label	{
//		self.Label	= other.Label
//	}
//	self.Username	= other.Username
//	self.Host		= other.Host
//	self.Port		= other.Port
//	self.Error		= other.Error
//
//	if	self.PrivateKeyPath	!= other.PrivateKeyPath	|| self.Password	!= other.Password	{
//		self.AuthMethods	= [] ssh.AuthMethod {}
//		self.ParsePrivateKey ( other.PrivateKeyPath )
//		self.SetPassword ( other.Password )
//	}
//
//	QueryIntervalSec
//	Commands
//	responses
//	error
//}

func ( self  * Servers )	Save ()	{
//	TODO
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