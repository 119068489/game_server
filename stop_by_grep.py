#!/usr/bin/python
# -*- coding: UTF-8 -*-
# CentOS 7 默认 Python 版本是 2.7
# CentOS 6 默认 Python 版本是 2.6
# 当前脚本文件要兼容 Python 2.6

import sys
import os
import readline # 使 Backspace 键和方向键不乱码

def main():
    if len(sys.argv) == 2:
        name  = sys.argv[1]
    else:
        name = raw_input("请输入要 kill 的进程名字(是模糊匹配,多组服务部署在一起时，注意别误杀了):")
        
    cmd = "ps -ef | grep %s | grep -v grep | awk '{print $2}' | xargs kill " % (name,)
    os.system(cmd)

if __name__ == "__main__":
    main()