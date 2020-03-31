// EOF requires the closing of the connection on the sending host i.e: RETR should close connection
// The server
// MUST close the data connection under the following conditions:
//    1. The server has completed sending data in a transfer mode
//       that requires a close to indicate EOF.
//    2. The server receives an ABORT command from the user.
//    3. The port specification is changed by a command from the
//       user.
//    4. The control connection is closed legally or otherwise.
//    5. An irrecoverable error condition occurs.
//
//In order to make FTP workable without needless error messages, the
// following minimum implementation is required for all servers:
//    TYPE - ASCII Non-print
//    MODE - Stream
//    STRUCTURE - File, Record
//    COMMANDS - USER, QUIT, PORT,
//               TYPE, MODE, STRU,
//                 for the default values
//               RETR, STOR,
//               NOOP.
// The default values for transfer parameters are:
//    TYPE - ASCII Non-print
//    MODE - Stream
//    STRU - File

package ftp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chibby0ne/go_book_exercises/chapter8/exercise8_2/flag"
	"github.com/chibby0ne/go_book_exercises/chapter8/exercise8_2/log"
)

type FtpConnection struct {
	// control connection
	//    The communication path between the USER-PI and SERVER-PI for
	//    the exchange of commands and replies.  This connection follows
	//    the Telnet Protocol.
	ControlConnnection net.Conn

	// data connection
	//    A full duplex connection over which data is transferred, in a
	//    specified mode and type. The data transferred may be a part of
	//    a file, an entire file or a number of files.  The path may be
	//    between a server-DTP and a user-DTP, or between two
	//    server-DTPs.
	DataConnection net.Conn

	// Defaults to 21. Port on which it is listening for ftp clients. Command line argument settable
	ControlPort int64

	// Defaults to 20. Port on which it will open a connection for data connection to ftp client. Command line argument settable
	DataPort int64

	// DataPort in client. Default is 20. Could be changed by PORT command
	DataPortClient string

	// DataHostClient. Default is same host as ControlConnection. Could be changed by PORT command
	DataHostClient string

	// Current user
	CurrentUser string

	// Binary type representation. Affects how files are copied
	Binary bool
}

func (f *FtpConnection) convertPortArgumentsToAddress(arg string) (string, string) {
	split := strings.Split(arg, ",")
	if len(split) != 6 {
		f.reply(SyntaxErrorInArgumentsReply)
		log.LogfVerbose("wrong number of arguments for convertion, format should be h1,h2,h3,h4,p1,p2: %s", split)
	}
	host := fmt.Sprintf("%s.%s.%s.%s", split[0], split[1], split[2], split[3])
	highBytePort, err := strconv.ParseInt(split[4], 10, 64)
	if err != nil {
		f.reply(SyntaxErrorInArgumentsReply)
		log.Fatalf("converting from int to string high byte: %s", err)
	}
	lowBytePort, err := strconv.ParseInt(split[5], 10, 64)
	if err != nil {
		f.reply(SyntaxErrorInArgumentsReply)
		log.Fatalf("converting from int to string low byte: %s", err)
	}
	port := fmt.Sprintf("%d", highBytePort<<8|lowBytePort)
	return host, port
}

// reply
//    A reply is an acknowledgment (positive or negative) sent from
//    server to user via the control connection in response to FTP
//    commands.  The general form of a reply is a completion code
//    (including error codes) followed by a text string.  The codes
//    are for use by programs and the text is usually intended for
//    human users.
type reply struct {
	completionCode int
	text           string
}

var (
	DataConnectionAlreadyOpenReply = reply{125, "Data connection already open; transfer starting"}
	FileStatusOKReply              = reply{150, "File status okay; about to open data connection."}

	CommandOKReply                       = reply{200, "Command okay"}
	AwaitingInputReply                   = reply{220, "Service ready for new user"} // used for greeting
	ServiceClosingControlConnectionReply = reply{221, "Service closing control connection"}
	CommandNotImplementedReply           = reply{202, "Command not implemented, superfluous at this site"}
	NameSystemTypeReply                  = reply{215, "UNIX system type"}
	DataConnectionOpenReply              = reply{225, "Data connection open; no transfer in progress"}
	ClosingDataConnection                = reply{226, "Closing data connection. Requested file action successful"}
	UserLoggedInReply                    = reply{230, "User logged in, proceed"}

	NeedAccountToLoginReply = reply{332, "Need account for login"}
	NeedPasswordReply       = reply{331, "User name okay, need password"}

	CannotOpenDataConnectionReply = reply{425, "Can't open data connection"}
	ConnectionClosedReply         = reply{426, "Connection closed; transfer aborted"}
	RequestFileNotTakenReply      = reply{450, "450 Requested file action not taken. File unavailable (e.g., file busy)"}
	RequestActionAbortedReply     = reply{451, "Requested action aborted: local error in processing."}

	SyntaxErrorCommandNotRecognizedReply   = reply{500, "Syntax error command not recognized"}
	SyntaxErrorInArgumentsReply            = reply{501, "Syntax error in parameters or arguments"}
	CommandNotImplementedForParameterReply = reply{504, "Command not implemented for that parameter"}
)

func (f *FtpConnection) reply(r reply) {
	log.LogfVerbose("reply: sending %v", r)
	fmt.Fprintf(f.ControlConnnection, "%s", r)
}

func (r reply) String() string {
	return fmt.Sprintf("%d %s \r\n", r.completionCode, r.text)
}

func getPort(addr net.Addr) (string, int64) {
	parts := strings.Split(addr.String(), ":")
	address, port := parts[0], parts[1]
	portNum, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return address, portNum
}

// Serve servers a single FTP connection to a FTP client
func (f *FtpConnection) Serve() {
	f.reply(AwaitingInputReply)
	scanner := bufio.NewScanner(f.ControlConnnection)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		cmd := strings.ToUpper(fields[0])
		var args []string
		if len(fields) > 0 {
			args = fields[1:]
		}
		switch cmd {
		case "USER":
			f.user(args)
		case "QUIT":
			f.quit(args)
		case "PORT":
			f.port(args)
		case "TYPE":
			f.typ(args)
		case "MODE":
			f.mode(args)
		case "STRU":
			f.stru(args)
		case "RETR":
			f.retr(args)
		case "STOR":
			f.stor(args)
		case "NOOP":
			f.noop(args)
		case "SYST":
			f.syst(args)
		default:
			f.reply(CommandNotImplementedReply)
			log.LogfVerbose("Command not yet defined: %q", cmd)
		}
	}
}

// USER NAME (USER)
//  The argument field is a Telnet string identifying the user.
//  The user identification is that which is required by the
//  server for access to its file system.  This command will
//  normally be the first command transmitted by the user after
//  the control connections are made (some servers may require
//  this).  Additional identification information in the form of
//  a password and/or an account command may also be required by
//  some servers.  Servers may allow a new USER command to be
//  entered at any point in order to change the access control
//  and/or accounting information.  This has the effect of
//  flushing any user, password, and account information already
//  supplied and beginning the login sequence again.  All
//  transfer parameters are unchanged and any file transfer in
//  progress is completed under the old access control
//  parameters.
func (f *FtpConnection) user(args []string) {
	if len(args) != 1 {
		f.reply(SyntaxErrorInArgumentsReply)
		log.LogfVerbose("wrong number of arguments supplied: %v", args)
		return
	}
	f.CurrentUser = args[0]
	log.LogfVerbose("user: logged in as :%v", args[0])
	f.reply(UserLoggedInReply)
}

// NOOP
// This command does not affect any parameters or previously
// entered commands. It specifies no action other than that the
// server send an OK reply.
func (f *FtpConnection) noop(args []string) {
	log.LogfVerbose("noop: nothing to be done")
	f.reply(CommandOKReply)
}

//LOGOUT (QUIT)
//  This command terminates a USER and if file transfer is not
//  in progress, the server closes the control connection.  If
//  file transfer is in progress, the connection will remain
//  open for result response and the server will then close it.
//  If the user-process is transferring files for several USERs
//  but does not wish to close and then reopen connections for
//  each, then the REIN command should be used instead of QUIT.
//  An unexpected close on the control connection will cause the
//  server to take the effective action of an abort (ABOR) and a
//  logout (QUIT).
func (f *FtpConnection) quit(args []string) {
	defer f.ControlConnnection.Close()
	f.reply(ServiceClosingControlConnectionReply)
}

// DATA PORT (PORT)
//   The argument is a HOST-PORT specification for the data port
//   to be used in data connection.  There are defaults for both
//   the user and server data ports, and under normal
//   circumstances this command and its reply are not needed.  If
//   this command is used, the argument is the concatenation of a
//   32-bit internet host address and a 16-bit TCP port address.
//   This address information is broken into 8-bit fields and the
//   value of each field is transmitted as a decimal number (in
//   character string representation).  The fields are separated
//   by commas.  A port command would be:
//      PORT h1,h2,h3,h4,p1,p2
//   where h1 is the high order 8 bits of the internet host
//   address.
func (f *FtpConnection) port(args []string) {
	if len(args) > 1 {
		f.reply(SyntaxErrorInArgumentsReply)
		log.LogfVerbose("wrong number of arguments supplied: %v", args)
		return
	}
	address, port := f.convertPortArgumentsToAddress(args[0])
	f.DataHostClient = address
	f.DataPortClient = port
	log.LogfVerbose("port: using %v:%v", address, port)
	f.reply(CommandOKReply)
}

// REPRESENTATION TYPE (TYPE)
//   The argument specifies the representation type as described
//   in the Section on Data Representation and Storage.  Several
//   types take a second parameter.  The first parameter is
//   denoted by a single Telnet character, as is the second
//   Format parameter for ASCII and EBCDIC; the second parameter
//   for local byte is a decimal integer to indicate Bytesize.
//   The parameters are separated by a <SP> (Space, ASCII code
//   32).

//   The following codes are assigned for type:

//                \    /
//      A - ASCII |    | N - Non-print
//                |-><-| T - Telnet format effectors
//      E - EBCDIC|    | C - Carriage Control (ASA)
//                /    \
//      I - Image

//      L <byte size> - Local byte Byte size
//      The default representation type is ASCII Non-print.  If the
//   Format parameter is changed, and later just the first
//   argument is changed, Format then returns to the Non-print
//   default.
func (f *FtpConnection) typ(args []string) {
	if len(args) == 0 || len(args) == 1 && strings.ToUpper(args[0]) != "I" || len(args) == 2 && strings.ToUpper(args[0]) == "I" {
		f.reply(SyntaxErrorInArgumentsReply)
		log.LogfVerbose("wrong number of arguments supplied: %v", args)
		return
	}
	log.LogfVerbose("type: args: %v", args)
	if args[0] == "A" && args[1] == "N" {
		f.Binary = false
	} else if args[0] == "I" || args[0] == "L" && args[1] == "8" {
		f.Binary = true
	} else {
		f.reply(CommandNotImplementedForParameterReply)
		return
	}
	f.reply(CommandOKReply)
}

// FILE STRUCTURE (STRU)
//    The argument is a single Telnet character code specifying
//    file structure described in the Section on Data
//    Representation and Storage.
//    The following codes are assigned for structure:
//       F - File (no record structure)
//       R - Record structure
//       P - Page structure
//    The default structure is File.
func (f *FtpConnection) stru(args []string) {
	if len(args) != 1 {
		f.reply(SyntaxErrorInArgumentsReply)
		log.LogfVerbose("wrong number of arguments supplied: %v", args)
		return
	}
	if args[0] == "F" {
		// Just implements default which is File
		log.LogfVerbose("stru: file structure")
		f.reply(CommandOKReply)
	} else {
		log.LogfVerbose("stru: not implemented for: %v", args)
		f.reply(CommandNotImplementedForParameterReply)
	}
}

//TRANSFER MODE (MODE)
// The argument is a single Telnet character code specifying
// the data transfer modes described in the Section on
// Transmission Modes.

// The following codes are assigned for transfer modes:

//    S - Stream
//    B - Block
//    C - Compressed

// The default transfer mode is Stream.
func (f *FtpConnection) mode(args []string) {
	if len(args) != 1 {
		f.reply(SyntaxErrorInArgumentsReply)
		log.LogfVerbose("wrong number of arguments supplied: %v", args)
		return
	}
	// Just implements default which is Stream
	if strings.ToUpper(args[0]) == "S" {
		log.LogfVerbose("mode: using stream")
		f.reply(CommandOKReply)
	} else {
		f.reply(CommandNotImplementedForParameterReply)
	}
}

// RETRIEVE (RETR)
//   This command causes the server-DTP to transfer a copy of the
//   file, specified in the pathname, to the server- or user-DTP
//   at the other end of the data connection.  The status and
//   contents of the file at the server site shall be unaffected.
func (f *FtpConnection) retr(args []string) {
	if len(args) != 1 {
		f.reply(SyntaxErrorInArgumentsReply)
		log.LogfVerbose("retr: wrong number of arguments supplied: %v", args)
		return
	}
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		log.LogfVerbose("retr: File: %v does not exist", args[0])
		f.reply(RequestFileNotTakenReply)
		return
	}
	file, err := os.Open(args[0])
	if err != nil {
		log.LogfVerbose("retr: error opening file: %v", err)
		f.reply(RequestFileNotTakenReply)
		return
	}
	log.LogfVerbose("retr: opened file: %v", file.Name())
	f.reply(FileStatusOKReply)
	conn, err := f.openDataConnection()
	if err != nil {
		log.LogfVerbose("retr: cannot open data connection: %s", err)
		f.reply(CannotOpenDataConnectionReply)
		return
	}
	defer conn.Close()
	_, err = io.Copy(conn, file)
	// bufferConn := bufio.NewWriter(conn)
	// _, err = bufferConn.ReadFrom(file)
	if err != nil {
		log.LogfVerbose("retr: Could not complete transfer of file: %s", err)
		f.reply(RequestActionAbortedReply)
	}
	log.LogfVerbose("retr: Copied file")
	f.reply(ClosingDataConnection)
	log.LogfVerbose("retr: Closing data connection")
}

// STORE (STOR)
//   This command causes the server-DTP to accept the data
//   transferred via the data connection and to store the data as
//   a file at the server site.  If the file specified in the
//   pathname exists at the server site, then its contents shall
//   be replaced by the data being transferred.  A new file is
//   created at the server site if the file specified in the
//   pathname does not already exist.
func (f *FtpConnection) stor(args []string) {
	if len(args) != 1 {
		f.reply(SyntaxErrorInArgumentsReply)
		log.LogfVerbose("stor: wrong number of arguments supplied: %v", args)
		return
	}
	log.LogfVerbose("stor: about to stor with args: %+v", args)
	dir, base := filepath.Dir(args[0]), filepath.Base(args[0])
	log.LogfVerbose("dir: %v, base: %v", dir, base)
	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		if errors.As(err, os.ErrPermission) {
			log.LogfVerbose("stor: Could not create directories due to permission error: %s")
		} else {
			log.LogfVerbose("stor: Could not create directories")
		}
		f.reply(RequestFileNotTakenReply)
		return
	}
	log.LogfVerbose("stor: created directory")
	file, err := os.Create(filepath.Join(dir, base))
	if err != nil {
		log.LogfVerbose("stor: error opening file: %v", err)
		f.reply(RequestFileNotTakenReply)
		return
	}
	log.LogfVerbose("stor: created file: %v", file.Name())
	f.reply(FileStatusOKReply)
	conn, err := f.openDataConnection()
	if err != nil {
		log.LogfVerbose("stor: cannot open data connection: %s", err)
		f.reply(CannotOpenDataConnectionReply)
		return
	}
	defer conn.Close()
	_, err = io.Copy(file, conn)
	if err != nil {
		log.LogfVerbose("stor: Could not complete transfer of file: %s", err)
		f.reply(RequestActionAbortedReply)
	}
	log.LogfVerbose("stor: Copied file")
	f.reply(ClosingDataConnection)
	log.LogfVerbose("stor: Closing data connection")
}

// SYSTEM (SYST)
//   This command is used to find out the type of operating
//   system at the server.  The reply shall have as its first
//   word one of the system names listed in the current version
//   of the Assigned Numbers document [4].
// inetutils' ftp uses this command after greeting
func (f *FtpConnection) syst(args []string) {
	if len(args) != 0 {
		f.reply(SyntaxErrorInArgumentsReply)
		log.LogfVerbose("wrong number of arguments supplied: %v", args)
		return
	}
	// In our case, it's only UNIX
	log.LogfVerbose("syst: replying with system type Unix")
	f.reply(NameSystemTypeReply)
}

func NewFTPConnection(conn net.Conn) *FtpConnection {
	ftpConnection := FtpConnection{
		ControlConnnection: conn,
		ControlPort:        *flag.Port,
		DataPort:           *flag.DataPort,
	}
	return &ftpConnection
}

// serveFTP serves concurrent FTP connections forever
func ServeFTP(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		go NewFTPConnection(conn).Serve()
	}
}

func (f *FtpConnection) openDataConnection() (net.Conn, error) {
	localAddress, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", f.DataPort))
	if err != nil {
		log.LogfVerbose("could not resolve local tcp address: %s", err)
		f.reply(CannotOpenDataConnectionReply)
		return nil, fmt.Errorf("could not resolve local tcp address: %s", err)
	}
	remoteAddress, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%s", f.DataHostClient, f.DataPortClient))
	if err != nil {
		log.LogfVerbose("could not resolve remote tcp address: %s", err)
		f.reply(CannotOpenDataConnectionReply)
		return nil, fmt.Errorf("could not resolve remote tcp address: %s", err)

	}
	conn, err := net.DialTCP("tcp4", localAddress, remoteAddress)
	if err != nil {
		log.Fatalf("could not establish connection: %s", err)
		f.reply(CannotOpenDataConnectionReply)
		return nil, fmt.Errorf("could not resolve remote tcp address: %s", err)
	}
	return conn, nil
}
