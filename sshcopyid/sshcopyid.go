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

package sshcopyid

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"

	sysadmSftp "github.com/wangyysde/sshclient/sftp"
	"github.com/wangyysde/sshclient/sshkey"
	"github.com/wangyysde/sysadm/utils"
)

var defaultIdentifyFile map[sshkey.Type]string = map[sshkey.Type]string {
	sshkey.KEY_RSA: ".ssh/id_rsa",
	sshkey.KEY_DSA: ".ssh/id_dsa",
	sshkey.KEY_ECDSA: ".ssh/id_ecdsa",
	sshkey.KEY_ED25519: ".ssh/id_ed25519",
}

func PubkeyAuthenticationTest(ip string, port int, user string, privateKey string ) (bool,error) {
	var idPath string = ""
	if strings.TrimSpace(privateKey) == "" {
		for _,v := range defaultIdentifyFile {
			idPath = filepath.Join(os.Getenv("HOME"),v)
			_,e := utils.CheckFileIsRead(idPath,"")
			if e == nil {
				break
			}
		}
	} else {
		idPath = strings.TrimSpace(privateKey) 
	}

	key, err := ioutil.ReadFile(idPath)
	if err != nil {
		return false, fmt.Errorf("unable to read private key: %v",err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	
	knowHostsPath := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	if ok,_ := utils.CheckFileExists(knowHostsPath,""); !ok {
		knowHostsFile, err := os.Create(knowHostsPath)
		if err != nil {
			return false, fmt.Errorf("can not create known_hosts file %s",err)
		}
		defer knowHostsFile.Close()
	}

	hostKeyCallback, err := kh.New(knowHostsPath)
	if err != nil {
		return false, fmt.Errorf("can not create hostkeycallback function %s", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	if err != nil {
		return false, fmt.Errorf("can not connect to server %s with port %d %s ", ip,port, err)

	} 
	client.Close()
	
	return true, nil 
}


func SshCopyId(ip string, port int, user string, password string, privateKey,pubKeyPath string) (bool, error){
	ok, _ := PubkeyAuthenticationTest(ip, port, user, privateKey)
	if ok {
		return true, nil
	}

	var pubKeyFile string = ""
	if strings.TrimSpace(pubKeyPath) == "" {
		for _,v := range defaultIdentifyFile {
			pubKeyFile = filepath.Join(os.Getenv("HOME"),v)
			pubKeyFile = pubKeyFile + ".pub"
			_,e := utils.CheckFileIsRead(pubKeyFile,"")
			if e == nil {
				break
			}
		}
	} else {
		pubKeyFile = strings.TrimSpace(pubKeyPath) 
	}

	// Read public key from file
	pubKey, err := os.ReadFile(pubKeyFile)
	if err != nil {
		return false, fmt.Errorf("Failed to read public key error: %v", err)
	}

	// open a new sftp connection
	client,err := sysadmSftp.ConnectSFTP(ip,user,password,"",port,false,true)
	if err != nil {
		return false,err
	}
	defer client.Close()

	// We expect working directory is SSH user's home directory
	workdir, err := client.Getwd()
	if err != nil {
		return false, fmt.Errorf("Failed to get working directory error: %v", err)
	}

	// Create .ssh and its parent folders if not exist
	remoteDir := filepath.Join(workdir, ".ssh")
	if err = client.MkdirAll(remoteDir); err != nil {
		return false, fmt.Errorf("Failed to create remote dir %s error: %v", remoteDir, err)
	}

	// Open or create authorized_keys to append public key
	remotePath := filepath.Join(remoteDir, "authorized_keys")
	file, err := client.OpenFile(remotePath, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	if err != nil {
		return false, fmt.Errorf("Failed to open file authorized_keys error: %v", err)
	}
	defer file.Close()
		
	if _, err = file.Write(pubKey); err != nil {
		return false, fmt.Errorf("Failed to write public key to remote authorized_keys error: %v", err)
	}

	return true,nil

}
