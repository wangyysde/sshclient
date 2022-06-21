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

#include <stdio.h>
#include <stdlib.h>
#include <termios.h>
#include <unistd.h>

#include "noecho.h"


/*
 * readpass switches the console to a non-echoing mode, reads a
 * line of standard input, and then switches the console back to
 * echoing mode.
 */
char *
readpass()
{
	struct termios	 term, restore;
	char		*password = NULL;
	size_t		 pw_size = 0;
	ssize_t		 pw_len;

	if (tcgetattr(STDIN_FILENO, &term) == -1)
		return NULL;

	restore = term;
	term.c_lflag &= ~ECHO;
	if (tcsetattr(STDIN_FILENO, TCSAFLUSH, &term) == -1)
		return NULL;

	pw_len = getline(&password, &pw_size, stdin);
	if (tcsetattr(STDIN_FILENO, TCSAFLUSH, &restore) == -1)
		return NULL;
	if (password != NULL)
		password[pw_len - 1] = (char)0;
	return password;
}
