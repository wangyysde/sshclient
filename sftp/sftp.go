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

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"github.com/wangyysde/sshclient/sshclient"
	"golang.org/x/crypto/ssh"
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

// puts a local file to remote server
// *sftp.Client a sftp connection opened by connectSFTP, srcFile resource file, dstFile 
// destination file.
// return error if error occured,otherwise return nil.
// Note: srcFile must be the path of a regular file.  
func Put(client *sftp.Client, srcFile ,dstFile string ) error {
	if client == nil {
		return fmt.Errorf("can not write file content on nil.")
	}

	if strings.TrimSpace(srcFile) == "" || strings.TrimSpace(dstFile) == "" {
		return fmt.Errorf("source file or destination file is empty.")
	}

	// Read content from source file
	srcContent, err := os.ReadFile(strings.TrimSpace(srcFile))
	if err != nil {
		return fmt.Errorf("Failed to read source file: %v", err)
	}
	
	srcInfo,e := os.Stat(strings.TrimSpace(srcFile))
	if e != nil {
		return  fmt.Errorf("get resouce file %s mode error %s",srcFile, e)
	}

	dstFp, err := client.OpenFile(strings.TrimSpace(dstFile), os.O_WRONLY|os.O_CREATE)
	if err != nil {
		return  fmt.Errorf("Failed to open destination file error: %v", err)
	}
	defer dstFp.Close()
	
	if _, err = dstFp.Write(srcContent); err != nil {
		return fmt.Errorf("write destionation file error: %v", err)
	}

	e = client.Chmod(strings.TrimSpace(dstFile), srcInfo.Mode())
    if e != nil {
		return fmt.Errorf("set remote dir %s mode error %s", dstFile,e)
    }	

	return nil 
}

// puts a local file or files in the srcFile to remote server
// *sftp.Client a sftp connection opened by connectSFTP, srcFile resource file or resource directory,dstFile
// destination file or directory. put will create a subdirectory named with the file name of srcFile in dstFile 
// if srcFile is a directory.
// return error if error occured,otherwise return nil.
func Mput(client *sftp.Client, srcFile, dstFile string,isEcho bool ) error {
	if client == nil {
		return fmt.Errorf("can not write file content on nil.")
	}

	if strings.TrimSpace(srcFile) == "" || strings.TrimSpace(dstFile) == "" {
		return fmt.Errorf("source file or destination file is empty.")
	}

	srcFile, err := filepath.Abs(srcFile)
	if err != nil {
		return fmt.Errorf("get absolute path error %s.",err)
	}

	fileInfo, err := os.Stat(srcFile)
	if err != nil {
		return fmt.Errorf("can not get file stats %s",err)
	}
	fileName := fileInfo.Name()
	if strings.TrimSpace(fileName) == "." || strings.TrimSpace(fileName) == ".." {
		return nil
	}

	if !fileInfo.IsDir() {
		e := Put(client,srcFile,dstFile)
		if isEcho {
			if e != nil {
				fmt.Printf("coping %s to %s error %s\n",srcFile, dstFile,e)
			} else {
				fmt.Printf("coping %s to %s successful\n",srcFile, dstFile)
			}
		}
	
		return e   
	}

	dir,err := os.ReadDir(srcFile)
	if err != nil {
		return fmt.Errorf("read the content in directory %s error %s", srcFile, err)	
	}
	
	for _,f := range dir {
        e :=  client.MkdirAll(dstFile)
		if e != nil {
			return fmt.Errorf("make dir %s on the remote server error %s", dstFile,e)
		}
	
		e = client.Chmod(dstFile, fileInfo.Mode()) 
        if e != nil {
			 return fmt.Errorf("set remote dir %s mode error %s", dstFile,e)
		}
		
		localSubFile := srcFile + string(os.PathSeparator) + f.Name()
		remoteSubFile := dstFile + string(os.PathSeparator) + f.Name()

		e = Mput(client,localSubFile, remoteSubFile,true)
		if e != nil {
			return fmt.Errorf("coping %s to %s on remote server error %s", localSubFile, remoteSubFile, e)
		}		
    }

	return nil
}
