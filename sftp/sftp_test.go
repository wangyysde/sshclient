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
	"os"
	"path/filepath"
	"testing"

	"github.com/wangyysde/sysadm/utils"
)

var (
	host string ="10.50.9.113"
	user string = "sysadm"
	password string = "Sysadm12345"
	privateKey string = ""
	port int = 22
	srcFile = "testdata/test_sftp.txt"
	dstFile = "/tmp/test_dst_sftp.txt"
	srcDir = "/tmp/testdata/testDirs"
	dstDir = "/tmp/testDirs"
)

func TestConnectWithPasswd(t *testing.T){
	client,err := ConnectSFTP(host,user,password,"",port,false,true)
	if err != nil {
		t.Errorf("connect to sftp server error: %s", err)
		return
	}
	defer client.Close()

	t.Log("connect to sftp server successful")

	// We expect working directory is SSH user's home directory
	workdir, err := client.Getwd()
	if err != nil {
		t.Errorf("get working dir  error: %s", err)
		return 
	}
	t.Logf("we got working directory is: %s",workdir)

	// Create .ssh and its parent folders if not exist
	remoteDir := filepath.Join(workdir, ".ssh")
	if err = client.MkdirAll(remoteDir); err != nil {
		t.Errorf("Failed to create remote dir %s error: %v", remoteDir, err)
		return 
	}

	// Open or create authorized_keys to append public key
	remotePath := filepath.Join(remoteDir, "authorized_keys")
	file, err := client.OpenFile(remotePath, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	if err != nil {
		t.Errorf("Failed to open file authorized_keys error: %v", err)
		return
	}
	defer file.Close()
	
	testContent := utils.Str2bytes("test message")
	if _, err = file.Write(testContent); err != nil {
		t.Errorf("Failed to write public key to remote authorized_keys error: %v", err)
	}

	return
}

func TestPut(t *testing.T){
	client,err := ConnectSFTP(host,user,password,"",port,false,true)
    if err != nil {
        t.Errorf("connect to sftp server error: %s", err)
        return
    }
    defer client.Close()

    t.Log("connect to sftp server successful")

    e := Put(client,srcFile,dstFile)
    if e != nil {
		t.Errorf("coping %s to %s on remote server error %s",srcFile,dstFile,e)
	} else {
		t.Logf("coping %s to %s on remote server successful",srcFile,dstFile)	
	}

	return

}

func TestMput(t *testing.T){
    client,err := ConnectSFTP(host,user,password,"",port,false,true)
    if err != nil {
        t.Errorf("connect to sftp server error: %s", err)
        return
    }
    defer client.Close()

    t.Log("connect to sftp server successful")

    e := Mput(client,srcDir,dstDir,true)
    if e != nil {
        t.Errorf("coping %s to %s on remote server error %s",srcDir,dstDir,e)
    } else {
        t.Logf("coping %s to %s on remote server successful",srcDir,dstDir)   
    }

    return

}

