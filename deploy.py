#!/usr/bin/python
# -*- coding:utf-8 -*-
# CentOS 7 默认 Python 版本是 2.7
# CentOS 6 默认 Python 版本是 2.6
# 当前脚本文件要兼容 Python 2.6

import os
import os.path
import sys
import readline # 使 Backspace 键和方向键不乱码

# CONFIG_NAME = "config_package"
SERVER_NAME  = "chat_server"
CHMOD_777_FILES = "%s start.py stop_by_pid_file.py stop_by_grep.py"%SERVER_NAME


def main():
	cwd = os.getcwd()
	cwd = os.path.basename(cwd)

	if len(sys.argv) != 2:
		while True:
			new_dir_name = raw_input("你要把当前目录名 {0} 改成什么新名字:".format(cwd))
			if new_dir_name != '':
				break
	else:
		new_dir_name = sys.argv[1]

	os.chdir("..")
	os.system("mv %s %s"%(cwd, new_dir_name))
	os.chdir(new_dir_name)

	print("给各个文件加权限...")
	os.system("chmod 777 {0}".format(CHMOD_777_FILES))

	print("建立各种软链接...")
	for name in os.listdir("./"):
		if os.path.isdir(name):
			os.chdir(name)
			os.system("ln -s ../%s %s"%(SERVER_NAME,SERVER_NAME))
			os.system("ln -s ../config_share.yaml config_share.yaml" )
			os.system("ln -s ../config_hall_secret.yaml config_hall_secret.yaml" )
			# os.system("ln -sfn ../%s/cheat cheat"%CONFIG_NAME)
			# os.system("ln -sfn ../%s/config config"%CONFIG_NAME)
			
			os.chdir("../")
	print ("部署成功。请 cd . 刷新当前目录就会看到 {0}。运行 python start.py 吧".format(new_dir_name))


if __name__ == "__main__":
	main()

