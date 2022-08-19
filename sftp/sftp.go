/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
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

// settings for sysadm system.
// settings in this file is temporary.they will be moved into Database
// TODO
// move the settings in this file into Database

package sftp

import(
	"fmt"
	"strings"
	"time"
	"io/ioutil"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"github.com/wangyysde/sshclient/sshclient"
)

// connectSFTP creates ssh config, tries to connect (dial) SSH server and
// creates new client with connection
func ConnectSFTP(host, username, password,privateKey string, port int, publicKeyAuth,ignoreHostKey bool ) (*sftp.Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	ignoreHostFunc := ssh.InsecureIgnoreHostKey()
	if !ignoreHostKey {
		hostFunc, err := sshclient.GetHostKeyCallbackFunc("")
		if err != nil {
			return nil, err
		}
		ignoreHostFunc = hostFunc
	}

	var config *ssh.ClientConfig
	if !publicKeyAuth {
		config = &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ignoreHostFunc,
			Timeout:         5 * time.Second,
		}
	} else {
		idPath := strings.TrimSpace(privateKey) 
		key, err := ioutil.ReadFile(idPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read private key: %v",err)
		}
		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(key)

		if err != nil {
			return nil, fmt.Errorf("can not parse private key: %v",err)
		}

		config = &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ignoreHostFunc,
			Timeout:         5 * time.Second,
		}
	}
	

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("Dial to SSH server %s error %s",host,err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new SFTP client %s",err)
	}

	return client, nil
}


