/*
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


Package sshkey provides Go handling of OpenSSH keys. It handles RSA
(protocol 2 only) and ECDSA keys, and aims to provide interoperability
between OpenSSH and Go programs.

The package can import public and private keys using the LoadPublicKey
and LoadPrivateKey functions; the LoadPublicKeyFile and LoadPrivateKeyfile
functions are wrappers around these functions to load the key from
a file. For example:

    // true tells LoadPublicKey to load the file locally; if false, it
    // will try to load the key over HTTP.
    pub, err := LoadPublicKeyFile("/home/user/.ssh/id_ecdsa.pub", true)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

In this example, the ECDSA key is in pub.Key. In order to be used in functions
that require a *ecdsa.PublicKey type, it must be typecast:

    ecpub := pub.Key.(*ecdsa.PublicKey)

The SSHPublicKey can be marshalled to OpenSSH format by using
MarshalPublicKey.

The package also provides support for generating new keys. The
GenerateSSHKey function can be used to generate a new key in the
appropriate Go package format (e.g. *ecdsa.PrivateKey). This key
can be marshalled into a PEM-encoded OpenSSH key using the
MarshalPrivate function.
*/
package sshkey
