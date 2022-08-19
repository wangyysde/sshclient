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

package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/wangyysde/sshclient/sshclient"
)

func main(){
	var (
		addr  string =  "172.28.2.4"
		port  string  = "22"
		user   string = "sysadm"
		passwd string = "Sysadm12345"
	//	prikey string = ""
	)

	remoteAddr := addr + ":" + port
	client, err := sshclient.DialWithPasswd(remoteAddr, user, passwd)
	if err != nil {
		fmt.Printf("DialWithPasswd err: %s", err)
		os.Exit(1)
	}
	defer client.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err = client.Cmd("/usr/bin/ls -l /").SetStdio(&stdout, &stderr).Run()
	if err != nil {
		fmt.Printf("Run command err: %s", err)
		os.Exit(2)
	}

	fmt.Printf("Content of stdout: %s\n",stdout.String())
	fmt.Printf("Content of stderr: %s\n",stderr.String())

	os.Exit(0)
}