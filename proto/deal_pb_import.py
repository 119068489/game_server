# -*- coding: utf-8 -*-
import re
import sys
import os

EXT_SRC = ".pb.go" # 原文件扩展名

def deal_a_proto_file(src):
	f1 = open(src, "r",encoding="utf8")
	txt = f1.read()
	f1.close()
	
	txt = txt.replace(r"/proto/share_message","/pb/share_message")
	txt = txt.replace(r"/proto/hall_game","/pb/hall_game")
	txt = txt.replace(r"/proto/client_game","/pb/client_game")
	txt = txt.replace(r"/proto/client_hall","/pb/client_hall")
	txt = txt.replace(r"/proto/client_server_proto","/pb/client_server")
	
	f2 = open(src, "w",encoding="utf8")
	
	f2.write(txt)
	f2.close()

def unpack_args(args):
	data = {}
	for arg in args:
		k,v = arg.split("=")
		data[k] = v
	return data

if __name__ == "__main__":
	args = unpack_args(sys.argv[1:])
	src_path = args["src"] # 源目录
	
	for src_name in os.listdir(src_path): # 不再使用 os.walk() 遍历，没有递归的需求
		if not src_name.endswith(EXT_SRC):
			continue
		srcFile = os.path.join(src_path, src_name)
		deal_a_proto_file(srcFile)


