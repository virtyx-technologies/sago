package main

import . "net"

type Node struct {
	Scheme          string  // TCP/IP scheme
	Host            string  // host name
	Ip              IP      // host ip address
	Port            int     // target port number
	Path            string  // 'path' portion of supplied URL, default '/'
	GetReq11        string  // http GET request for the URL
	startTlsOpts    string  // '--starttls' options for openssl, customised according to the protocol
	features        map[string]bool // Features discovered at run time

	OptimalProtocol string  // not sure what this is, might be redundant
}

func (n *Node) hostPort() string{
	return JoinHostPort(n.Ip.String(), string(n.Port))
}

func (n *Node) IsIpv4Addr() bool {
	return n.Ip.To4()  != nil
}

func (n *Node) IsIpv6Addr() bool {
	return n.Ip.To16()  != nil
}

