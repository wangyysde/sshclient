/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2021 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
* @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html
 */

// Package sshclient implements an SSH client.

package sshclient

import (
	"bytes"
	"io"

	"golang.org/x/crypto/ssh"
)

type remoteScriptType byte
type remoteShellType byte

const (
	cmdLine remoteScriptType = iota
	rawScript
	scriptFile

	interactiveShell remoteShellType = iota
	nonInteractiveShell
)

// A Client implements an SSH client that supports running commands and scripts remotely.
type Client struct {
	client *ssh.Client
}

// A RemoteScript represents script that can be run remotely.
type RemoteScript struct {
	client     *ssh.Client
	_type      remoteScriptType
	script     *bytes.Buffer
	scriptFile string
	err        error

	stdout io.Writer
	stderr io.Writer
}

// A RemoteShell represents a login shell on the client.
type RemoteShell struct {
	client         *ssh.Client
	requestPty     bool
	terminalConfig *TerminalConfig

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// A TerminalConfig represents the configuration for an interactive shell session.
type TerminalConfig struct {
	Term   string
	Height int
	Weight int
	Modes  ssh.TerminalModes
}