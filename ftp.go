package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

// Status defines FTP status code.
const (
	StatusCommandOK              = 200
	StatusServiceReadyForNewUser = 220
	StatusEnteringPassiveMode    = 227
	StatusUserLoggedIn           = 230
	StatusNeedPassword           = 331
)

// FTPClient represents a FTP client.
type FTPClient struct {
	host     string
	port     int
	conn     *textproto.Conn
	logger   *log.Logger
	dataHost string
	dataPort int
}

// NewFTPClient constructs a new FTP client.
func NewFTPClient(host string, port int) (*FTPClient, error) {
	return &FTPClient{
		host:   host,
		port:   port,
		logger: log.New(ioutil.Discard, "", 0),
	}, nil
}

// Close closes FTP connection.
func (c *FTPClient) Close() {
	if c.conn != nil {
		if _, _, err := c.Cmd("QUIT"); err != nil {
			c.logger.Printf("failed to send QUIT command: %v", err)
		}
		err := c.conn.Close()
		if err != nil {
			c.logger.Printf("failed to close control connection: %v", err)
		}
	}
}

// SetLogger updates logger. Default is discard.
func (c *FTPClient) SetLogger(logger *log.Logger) {
	c.logger = logger
}

// Login logins to the SSH server.
func (c *FTPClient) Login(user, password string) error {
	if err := c.open(); err != nil {
		return err
	}
	if code, msg, err := c.Cmd("USER " + user); err != nil {
		return fmt.Errorf("failed to send USER command: %w", err)
	} else if code != StatusNeedPassword {
		return fmt.Errorf("failed to execute USER command: FTP %d %s", code, msg)
	}
	if code, msg, err := c.Cmd("PASS " + password); err != nil {
		return fmt.Errorf("failed to send PASS command: %w", err)
	} else if code != StatusUserLoggedIn {
		return fmt.Errorf("failed to execute PASS command: FTP %d %s", code, msg)
	}
	if code, msg, err := c.Cmd("TYPE I"); err != nil {
		return fmt.Errorf("failed to send TYPE command: %w", err)
	} else if code != StatusCommandOK {
		return fmt.Errorf("failed to execute TYPE command: FTP %d %s", code, msg)
	}
	return nil
}

// Cmd executes a FTP command.
func (c *FTPClient) Cmd(text string) (code int, msg string, err error) {
	c.logger.Print("> " + text)
	if _, err = c.conn.Cmd(text); err != nil {
		return
	}
	code, msg, err = c.conn.ReadResponse(0)
	if err != nil {
		c.logger.Print(err)
		return
	}
	c.logger.Printf("< %d %s", code, msg)
	return
}

// DataCmd executes a FTP command using data connection.
func (c *FTPClient) DataCmd(text string) (string, error) {
	addr, err := c.passiveMode()
	if err != nil {
		return "", err
	}
	c.logger.Printf("> %s", addr)
	dataConn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("failed to dial data connection: %w", err)
	}
	defer func() {
		if err := dataConn.Close(); err != nil {
			c.logger.Printf("failed to close data connection: %v", err)
		}
		code, msg, err := c.conn.ReadResponse(0)
		if err != nil {
			c.logger.Printf("transfer incomplete: %v", err)
		}
		c.logger.Printf("< %d %s", code, msg)
	}()
	if _, _, err := c.Cmd(text); err != nil {
		return "", err
	}
	resp, err := ioutil.ReadAll(dataConn)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	return string(resp), nil
}

func (c *FTPClient) open() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("failed to close previous connection: %w", err)
		}
	}
	c.logger.Printf("> %s:%d", c.host, c.port)
	conn, err := textproto.Dial("tcp", fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		c.logger.Print(err)
		return fmt.Errorf("failed to dial %s:%d: %w", c.host, c.port, err)
	}
	code, msg, err := conn.ReadResponse(StatusServiceReadyForNewUser)
	if err != nil {
		c.logger.Print(err)
		return fmt.Errorf("server is not ready: FTP %d %s", code, msg)
	}
	c.logger.Printf("< %d %s", code, msg)
	c.conn = conn
	return nil
}

// passiveMode enables data connection using passive mode.
func (c *FTPClient) passiveMode() (string, error) {
	code, msg, err := c.Cmd("PASV")
	if err != nil {
		return "", fmt.Errorf("failed to send PASV command: %w", err)
	}
	if code != StatusEnteringPassiveMode {
		return "", fmt.Errorf("failed to execute PASV command: FTP %d %s", code, msg)
	}
	hp := strings.Split(msg[strings.LastIndex(msg, "(")+1:strings.LastIndex(msg, ")")], ",")
	if len(hp) != 6 {
		return "", fmt.Errorf("unexpected PASV response: %s", msg)
	}
	p1, err := strconv.Atoi(hp[4])
	if err != nil {
		return "", fmt.Errorf("unexpected PASV port response: %s", msg)
	}
	p2, err := strconv.Atoi(hp[5])
	if err != nil {
		return "", fmt.Errorf("unexpected PASV port response: %s", msg)
	}
	return net.JoinHostPort(strings.Join(hp[0:4], "."), strconv.Itoa(p1*256+p2)), nil
}
