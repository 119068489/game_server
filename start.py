#!/usr/bin/python
# -*- coding: UTF-8 -*-
# CentOS 7 默认 Python 版本是 2.7
# CentOS 6 默认 Python 版本是 2.6
# 当前脚本文件要兼容 Python 2.6

import sys
import shutil
import time
import os
import os.path
import readline # 使 Backspace 键和方向键不乱码

EXECUTE_NAME = "chat_server"

# 约定: 可执行目录必须以游戏名相同，或者游戏名作前缀并以-分隔
# 比如： fish fish-2 fish-3 fish-new 都会当成是 fish 游戏
# 或者： hall hall-2 hall-3 hall-new 都会当成是 hall

EXE_DIR_LIST = ["login","hall", "game", "backstage","shop","statistics","square"] # 不包含 sub_game 这个子类型

PID_FILE_NAME = "pid"

def get_dict_from_pid_file(): # 在 stop_by_pid_file.py 也存在此函数定义
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


# 是否合法的可执行目录名称
def is_exe_dir(dir_name):
    name = dir_name.split("-")[0] #
    return name in EXE_DIR_LIST

# 从磁盘中取得可执行目录列表
def get_real_exe_dir_list():
    real_exe_dir_list = os.listdir('./')
    for i in range(len(real_exe_dir_list)-1, -1, -1):# 排除文件，只剩目录。在 for 循环中逆序删，避免跳跃
        path = real_exe_dir_list[i]
        if not os.path.isdir(path):
            real_exe_dir_list.pop(i)
            continue

        if not is_exe_dir(path):
            real_exe_dir_list.pop(i)
            continue
    return real_exe_dir_list

def main():
    os.system("ulimit -u 8192")
    os.system("ulimit -n 8192")
    real_exe_dir_list = get_real_exe_dir_list()
    need_input = True
    if len(sys.argv) > 1:
        dir_name = sys.argv[1]
        if dir_name == "all" or dir_name in real_exe_dir_list:
            need_input = False
        else:
            print("你输入的 {0} 不是有效的 exe 目录".format(dir_name))

    while need_input:
        dir_name = raw_input("请输入你想要启动的目录名.可以是这些 {0} 之一 (all 表示启动所有子目录): ".format(real_exe_dir_list))
        if dir_name == '':
            continue
        if dir_name == "all" or dir_name in real_exe_dir_list:
            break

    if len(sys.argv) > 2:
        group = sys.argv[2]
    else: # 服务器组标识
        cwd = os.getcwd()
        group = os.path.basename(cwd)
        s = raw_input("请输入服务器组标识。允许不输入,直接回车则是当前目录名 {0} : ".format(group))
        if s != "":
            group = s

    if dir_name == "all":
        if os.path.exists(PID_FILE_NAME):
           os.system("rm -rf %s"%PID_FILE_NAME)

        for exe_dir in real_exe_dir_list:
            os.chdir("./%s"%exe_dir)
            run_execute(exe_dir, group)
            os.chdir("..")
            time.sleep(1)
    else:
        if not os.path.exists(PID_FILE_NAME): # 生成 pid 文件
            os.system("touch %s"%PID_FILE_NAME)

        pid_info = get_dict_from_pid_file()
        if dir_name in pid_info:
            raise Exception("已经存在键值为 %s 的 pid,该目录已启动了？！" % dir_name)

        os.chdir(dir_name)
        run_execute(dir_name, group)

def run_execute(exe_dir, group):
    time_str = time.strftime("%Y-%m-%d_%H_%M_%S%z", time.localtime())
    cmd = 'nohup ./{exe} -app={app} [{exe_dir}@{group}] > /dev/null 2> ./logs/{app}_std_error.log &echo "{exe_dir}:$!" >> ../{pid}'
    cmd = cmd.format(exe=EXECUTE_NAME, app=exe_dir, group=group,exe_dir=exe_dir, time=time_str, pid=PID_FILE_NAME)

    os.system(cmd)

if __name__ == "__main__":
    main()
