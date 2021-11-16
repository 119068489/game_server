# -*- coding: UTF-8 -*-

# 当前脚本只能在 Windows 下运行
# 可以使用 Python 3 的语法，Windows 下大家约定用 Python 3.7 以上

import os
import shutil
import subprocess
import sys
import hashlib
import socket
import time
import build_linux_yaml_template as template
NOW_TEXT = time.strftime("%Y-%m-%d_%H_%M_%S%z", time.localtime())
# 要保证有 2 位数来表示子游戏类型，只能 2 位，不能多不能少。
DIR_LIST = ["hall","login","backstage","shop","statistics"]
DIR_NAME = "linux_server_" + NOW_TEXT# 临时打包目录，打包完后删除
EXECUTE_NAME = "chat_server.exe" # 打包出来的可执行文件名

NEED_HASH_DIR_LIST = [ # 后续优化: 不手工维护这个表了，自动递归全部目录就好了，发现 *.go 就 hash
    "execute",
    "hall",
    "login",
    "shop",
    "backstage",
    "mongo_init",
    "easygo",
    "for_game",
    "pb",
    "deleter",
    ]

FINGERPRINT_FILE_NAME = DIR_NAME + "/fingerprint_{now}.txt".format(now=NOW_TEXT)
# CONFIG_DIR_NAME = "config_package"
DEFAULT_HOST = '192.168.50.27'
DEFAULT_DB='192.168.50.27'

# 检查有没有改了没有提交的或是新的 *.go 文件
def has_change_or_new_file():
    code, s = subprocess.getstatusoutput("git status")
    if code != 0:
        raise Exception(s)
    return ".go" in s

def build(dir_name):
    os.chdir("execute")
    print("准备编译可执行文件 {} ...".format(EXECUTE_NAME))

    #text = "set CGO_ENABLED=0&&set GOOS=linux&&set GOARCH=amd64&&go build -o ../{dir_name}/{exe_name}".format(dir_name=dir_name, exe_name=EXECUTE_NAME)
    text = "go build -o ../{dir_name}/{exe_name}".format(dir_name=dir_name, exe_name=EXECUTE_NAME)
    code,s = subprocess.getstatusoutput(text)  #必须在同一个线程运行 否则不成功
    if code != 0:
        raise Exception(s)
    print("编译成功")
    os.chdir(".." )

def deal_yaml_json_py_etc(dir_name, dir_list, is_full, group, protocol, host,db):       # 打包linux版本文件夹
    if not is_full:
        return
    os.chdir("execute")

    os.chdir("../%s"%dir_name)
    # 复制文件到当前目录
    shutil.copy("../start.py", "./")
    shutil.copy("../stop_by_pid_file.py", "./")
    shutil.copy("../stop_by_grep.py", "./")

    shutil.copy("../backup.py", "./")
    shutil.copy("../deploy.py", "./")
    
    content = template.TEMPLATE_SHARE.format(group=group,center="127.0.0.1",db="127.0.0.1")
    with open("./config_share.yaml", "w", encoding="utf-8") as f:
        f.write(content)
    content = template.TEMPLATE_HALL_SECRET # 直接写，无需 format .format(group=group, host=host)
    with open("./config_hall_secret.yaml", "w", encoding="utf-8") as f:
        f.write(content)
    # os.mkdir(CONFIG_DIR_NAME)
    # os.chdir(CONFIG_DIR_NAME)
    # os.system('xcopy "../../cheat" "cheat" /s /e  /i /y')
    # os.system('xcopy "../../config" "config" /s /e  /i /y')
    # os.chdir("../")

    for dir in dir_list:     #把配置文件复制到各个文件夹下
        os.mkdir(dir)
        print("创建 %s\t子目录,并生成了 yaml 配置文件进去 "%(dir,))
        os.chdir(dir)
        if dir == "hall":
            content = template.TEMPLATE_HALL.format(group=group, host=host)
            with open("./config_hall.yaml", "w", encoding="utf-8") as f:
                f.write(content)
        elif dir == "login":
            content = template.TEMPLATE_LOGIN.format(group=group, host=host)
            with open("./config_login.yaml", "w", encoding="utf-8") as f:
                f.write(content)
        elif dir == "shop":
            content = template.TEMPLATE_SHOP.format(group=group, host=host)
            with open("./config_shop.yaml", "w", encoding="utf-8") as f:
                f.write(content)
        elif dir == "backstage":
            content = template.TEMPLATE_BACKSTAGE.format(group=group, host=host)
            with open("./config_backstage.yaml", "w", encoding="utf-8") as f:
                f.write(content)
            shutil.copy("../../backstage/version.json", "./")
            shutil.copy("../../backstage/tfserver.json", "./")
        elif dir == "statistics":
            content = template.TEMPLATE_STATISTICS.format(group=group, host=host)
            with open("./config_statistics.yaml", "w", encoding="utf-8") as f:
                f.write(content)
        else:
            raise Exception("未知的目录 "+ dir)
        os.mkdir("logs")
        os.chdir("../")

    os.chdir("../")

def package_zip(dir_name, is_full): # 把打包文件夹压缩成zip文件
    print("开始压缩 %s 目录,耗时较长，耐心等候 ...... " %(dir_name,))

    if is_full:
        t = "full"
    else:
        t = "execute"
    name = "%s_%s.zip" %(dir_name, t)
    text = "7z.exe -tZip a %s ./%s -mx9"%(name, dir_name)
    code, s = subprocess.getstatusoutput(text)
    if code != 0:
        text = "安装7z压缩软件了吗？？？设置7z的环境变量了吗？？？"
        raise Exception(text + s)
    print("压缩 OK,包名是 "+name)

def remove_dir(dir_name):    # 删除打包文件夹
    if os.path.exists(dir_name):
        print("删除临时打包目录 "+ dir_name)
        shutil.rmtree(dir_name)
        

def hash_file(file_name): # hash 出 md5 值
    if not os.path.isfile(file_name):
        return
    myhash = hashlib.md5()
    with open(file_name,'rb') as f:
        while True:
            b = f.read(8096)
            if not b:
                break
            myhash.update(b)
    return myhash.hexdigest()

def hash_all_file(dir_name): # 获取到所有当前路径下的文件
    lst = []
    for (root, dirs, files) in os.walk(dir_name):
        _ = dirs
        for file_name in files:
            s1 = hash_file(root+"\\"+file_name)
            s2 = "%s\\%s: %s\n" % (root,file_name, s1)
            
            lst.append(s2)
    return "".join(lst)

def gen_fingerprint_file(fingerprint_file, need_hash_dir_list, branch_name):  # 哈希 *.go 代码文件
    if os.path.exists(fingerprint_file):     # 检测如果有这个文件就删除新建
        os.remove(fingerprint_file)

    with open(fingerprint_file,"a",encoding="utf8") as f:
        host_name = socket.gethostname()  # 获取本机计算机名
        f.write("计算机名: %s\n"%host_name)
        f.write("打包时间: %s\n" % NOW_TEXT)
        f.write("打包工作目录: %s\n" % os.getcwd())
        f.write("打包分支名: {}\n".format(branch_name))

        # 获取当前提交版本号
        code, s = subprocess.getstatusoutput("git rev-parse HEAD")
        f.write("最后 Commit: %s\n" % s)
        if code != 0:
            raise Exception(s)

        # 获取当前环境 Golang 版本
        code,s = subprocess.getstatusoutput("go version")
        if code != 0:
            raise Exception(s)

        f.write("打包机器 Golang 版本: %s" % s)
        f.write("\n")

        digest = hash_file("./{dir_name}/{exe_name}".format(dir_name=DIR_NAME, exe_name=EXECUTE_NAME))
        f.write("可执行文件 {} MD5 值: {}\n".format(EXECUTE_NAME, digest))

        f.write("\n各源代码文件 MD5 值:\n")
        for dir_name in need_hash_dir_list:     # 循环遍历所有需要 hash 的目录
            text = hash_all_file(dir_name)
            f.write(text)
    
    print("生成各 *.go 源码文件的 hash 值成功")

def main():
    code, branch_name = subprocess.getstatusoutput("git symbolic-ref --short -q HEAD")
    if code != 0:
        raise Exception(branch_name)
    if branch_name != "master":
        while True:
            q = input("严重警告!!!!!! 当前分支是 {}，你真的要对这个分支而不是 master 进行打包 (输入 y 或 n): ".format(branch_name))
            if q == "":
               continue
            elif q == 'y':
                break
            else:
                print("中止打包")
                return

    if has_change_or_new_file():
        while True:
            q = input("严重警告!!!!!! 发现有新的或是改动未提交的 go 文件，是否仍要继续打包? (输入 y 或 n): ")
            if q == "":
               continue
            elif q == 'y':
                break
            else:
                print("中止打包")
                return


    while True:
        s = input("打完整包还是只打可执行文件?(输入 full 代表打完整包，输入 exe 代表打可执行文件): ")
        if s == "":
            continue
        if s in ["full", "exe"]:
            is_full = {"full":True, "exe":False}[s]
            break

    if is_full:
        while True:
            group = input("请输入服务器组,用于各监听端口的最后一位数,有效值为 0 - 9: ")
            if len(group) == 1 and group.isdigit():
                break
        while True:
            protocol = input("游戏客户端和服务器走什么协议？请输入 ws 或 wss : ")
            if protocol in ("ws", "wss"):
                break

        host = input("请输入目标服务器的外网 IP 或域名(直接回车则是 {}): ".format(DEFAULT_HOST))
        if host == "":
            host = DEFAULT_HOST
        db = input("请输入mongodb的IP(直接回车则是 {}): ".format(DEFAULT_DB))
        if db == "":
            db = DEFAULT_DB
        while True:
            is_all = input("打包服务器all表示全部[login 、hall、backstage、shop、statistics]其中一个): ")
            if is_all == "all" or is_all in DIR_LIST:
                break

    while True:
        s = input("是否压缩? (输入 y 或 n): ")
        if s == "":
            continue
        if s in ["y", "n"]:
            compress = {"y":True, "n":False}[s]
            break
    remove_dir(DIR_NAME)
    os.mkdir(DIR_NAME)

    build(DIR_NAME)
    gen_fingerprint_file(FINGERPRINT_FILE_NAME, NEED_HASH_DIR_LIST, branch_name)

    if is_full:
        server_list = []
        if is_all =="all":
            server_list=DIR_LIST
        else:
            server_list=[is_all]
        deal_yaml_json_py_etc(DIR_NAME, server_list, is_full, group, protocol, host,db)

    if compress:
        package_zip(DIR_NAME, is_full) # 压缩
        remove_dir(DIR_NAME) # 删除临时打包文件夹

if __name__ == "__main__":
    main()