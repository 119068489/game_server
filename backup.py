#!/usr/bin/python
# -*- coding:utf-8 -*-

# CentOS 7 默认 Python 版本是 2.7
# CentOS 6 默认 Python 版本是 2.6
# 当前脚本文件要兼容 Python 2.6

import os
import time
import sys
import readline # 使 Backspace 键和方向键不乱码

BACKUP_DIR = "server_backup" # 专门存储目录

def main():     # 备份功能
    cwd = os.getcwd()
    need_backup = os.path.basename(cwd)
    os.chdir("../")

    if not os.path.exists(BACKUP_DIR):
        os.system("mkdir %s" % BACKUP_DIR)


    
    new_name = "%s_bak_%s" % (need_backup, time.strftime("%Y_%m_%d_%H_%M%z", time.localtime()))
    os.system("mv  %s  %s" % (need_backup, new_name))
    # -y: 遇到软链接，复制软链接，不复制原文件
    os.system("zip -ry %s.zip %s " % (new_name, new_name))
    os.system("mv %s.zip ./%s/" % (new_name, BACKUP_DIR))
    os.system("rm -rf %s" % new_name)


if __name__ == "__main__":
	main()