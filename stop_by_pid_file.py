#!/usr/bin/python
# -*- coding: UTF-8 -*-
# CentOS 7 默认 Python 版本是 2.7
# CentOS 6 默认 Python 版本是 2.6
# 当前脚本文件要兼容 Python 2.6

# 读取 pid 文件中记录的进程 id 杀进程
import sys
import os
import readline # 使 Backspace 键和方向键不乱码

PID_FILE_NAME = "pid"

def get_dict_from_pid_file(): # 在 start.py 也存在此函数定义
    pid_info = {}
    if not os.path.exists(PID_FILE_NAME):
        return pid_info

    with open('./%s' % PID_FILE_NAME, 'r') as f:
        s = f.read()
    if s == "":
        return pid_info

    for line in s.split("\n"):
        if line == "":
            continue
        exe_dir, pid = line.split(":")
        pid_info[exe_dir] = int(pid)
    return pid_info


def write_pid_file(pid_info):
    with open('./%s' % PID_FILE_NAME, 'w') as f:
       for exe_dir, pid in pid_info.iteritems():
           f.write("%s:%s\n"%(exe_dir, pid))

def main():
    pid_file = PID_FILE_NAME
    if not os.path.exists(pid_file):
        print("没有 pid 文件存在.无能为力")
        return
    pid_info = get_dict_from_pid_file()

    need_input = True
    if len(sys.argv) > 1:
        dir_name = sys.argv[1]
        if dir_name == "all" or dir_name in pid_info:
            need_input = False
        else:
            print("你输入的目录 {0} 在 pid 文件中找不到".format(dir_name))

    while need_input:
        prompt = "请输入你想要 kill 的目录名. pid 文件中记录的有: {0}  (all 表示 kill 全部): ".format(pid_info.keys())
        dir_name = raw_input(prompt)
        if dir_name == '':
            continue
        if dir_name != "all" and dir_name not in pid_info:
            print("目录名 {0} 在 pid 文件中不存在".format(dir_name))
            continue
        break

    if dir_name == "all": #关闭所有子游戏进程
        for pid in pid_info.values():
            os.system("kill  %d"%pid)
        os.system("rm -rf %s" % pid_file)
        
    else: # 关闭某个子游戏进程
        pid = pid_info[dir_name]
        os.system("kill  %d" % pid)
        del pid_info[dir_name]
        write_pid_file(pid_info) # 覆盖 pid 文件

if __name__ == "__main__":
    main()