package gremtune

import "time"

//DialerConfig is the struct for defining configuration for WebSocket dialer
type DialerConfig func(*Websocket)

//SetAuthentication sets on dialer credentials for authentication
func SetAuthentication(username string, password string) DialerConfig {
	return func(c *Websocket) {
		c.auth = &auth{username: username, password: password}
	}
}

//SetTimeout sets the dial handshake timeout
func SetTimeout(seconds int) DialerConfig {
	return func(c *Websocket) {
		c.timeout = time.Duration(seconds) * time.Second
	}
}

//SetPingInterval sets the interval of ping sending for know is
//connection is alive and in consequence the client is connected
func SetPingInterval(seconds int) DialerConfig {
	return func(c *Websocket) {
		c.pingInterval = time.Duration(seconds) * time.Second
	}
}

//SetWritingWait sets the time for waiting that writing occur
func SetWritingWait(seconds int) DialerConfig {
	return func(c *Websocket) {
		c.writingWait = time.Duration(seconds) * time.Second
	}
}

//SetReadingWait sets the time for waiting that reading occur
func SetReadingWait(seconds int) DialerConfig {
	return func(c *Websocket) {
		c.readingWait = time.Duration(seconds) * time.Second
	}
}

//SetBufferSize sets the read/write buffer size
func SetBufferSize(readBufferSize int, writeBufferSize int) DialerConfig {
	return func(c *Websocket) {
		c.readBufSize = readBufferSize
		c.writeBufSize = writeBufferSize
	}
}

// websocketDialerFactoryFun exchange/ set the factory function used to create the dialer which
// is then used to open the websocket connection.
// This function is not exported on purpose, it should only used for injection and mocking in tests!!
func websocketDialerFactoryFun(wsDialerFactory websocketDialerFactory) DialerConfig {
	return func(c *Websocket) {
		c.wsDialerFactory = wsDialerFactory
	}
}
