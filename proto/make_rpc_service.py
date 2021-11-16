# -*- coding: utf-8 -*-
import re
import sys
import os

EXT_SRC = ".proto" # 原文件扩展名
EXT_DEST = ".go" # 目标文件扩展名


def deal_a_proto_file(src, dest):
	f = open(src, "rt",encoding="utf8")
	txt = f.read()
	f.close()
	
	package_name = parse_package_name(txt)

	rpc_codes = []
	service_names = []
	# 一般就是处理 2 个 service
	i = 0
	up, down = [], []
	import_name_set = set()
	for service_name,lines in re.findall(r"\nservice\s(\w+)\s?\{([\s\S]+?)\}", txt):
		i += 1
		methods, interface_items = [], []
		for line in lines.split("\n"):
			a_map_item,a_rpc_method, a_interface_item, package_set, ok = deal_a_rpc(package_name, line, service_name)
			if ok:
				if i == 1:
					up.append(a_map_item)
				else:
					down.append(a_map_item)
				methods.append(a_rpc_method)
				interface_items.append(a_interface_item)

				import_name_set = import_name_set|package_set
		rpc_methods = "\n".join(methods)
		rpc_interfaces = "\n".join(interface_items)

		a_service_code = A_SERVICE_CODE.format(service_name=service_name,rpc_methods=rpc_methods,rpc_interfaces=rpc_interfaces)
		rpc_codes.append(a_service_code)
		
		service_names.append(service_name)
		
	if rpc_codes:
		service_codes = "\n\n// ==========================================================\n".join(rpc_codes)
		import_name  = ""
		for name in import_name_set:
			import_name +=  '"game_server/pb/{}"\n'.format(name)
		schema = SCHEMA_CODE.format(package_name=package_name,service_codes=service_codes,import_name=import_name)
		f = open(dest, "w",encoding="utf-8")
		f.write(schema)
		f.close()
	return service_names, up, down

def parse_package_name(s):
	m = re.match(r"package (\S+);", s)
	if m:
		name = m.group(1)
		return name
	return ""


def deal_a_rpc(package_name, line,service_name):
	line = line.strip()
	m = re.match(r"rpc.(\S+)\((\S+)\)\s?returns\((\S+)\);", line)
	package_set = set()
	if m:
		method_name = m.group(1)
		request_msg = m.group(2)
		response_msg = m.group(3)
		a_map_item = '\t"%s": easygo.Pair{"%s", "%s"},' % (method_name, request_msg, response_msg)
		req_msg_package, msg = request_msg.split(".") # 同一个包下的 msg 要去掉包前缀
		if req_msg_package == package_name:
			request_msg = msg
		else:
			if req_msg_package != "base" and req_msg_package not in package_set:
				package_set.add(req_msg_package)

		resp_msg_package, msg = response_msg.split(".") # 同一个包下的 msg 要去掉包前缀
		if resp_msg_package == package_name:
			response_msg = msg
		else:
			if resp_msg_package != "base" and resp_msg_package not in package_set:
				package_set.add(resp_msg_package)

		if response_msg!="base.NoReturn":
			method_templ = A_RPC_METHOD_1
			interface_templ = A_INTERFACE_ITEM_1
		else:
			method_templ = A_RPC_METHOD_2
			interface_templ = A_INTERFACE_ITEM_2

		a_rpc_method = method_templ.format(service_name=service_name,method_name=method_name,request_msg=request_msg,response_msg=response_msg)
		a_interface_item = interface_templ.format(service_name=service_name,method_name=method_name,request_msg=request_msg,response_msg=response_msg)
		return a_map_item,a_rpc_method, a_interface_item,package_set, True
	return "","","",package_set,False

def unpack_args(args):
	data = {}
	for arg in args:
		k,v = arg.split("=")
		data[k] = v
	return data
#----------------------------
SCHEMA_MAP = """package {rpc_map_file}

import (
	"game_server/easygo"
)
var UpRpc = map[string]easygo.Pair{{
{up_rpc_items}
}}

var DownRpc = map[string]easygo.Pair{{
{down_rpc_items}
}}"""



SCHEMA_CODE = """package {package_name}

import (
	"game_server/easygo"
	"game_server/easygo/base"
	{import_name}
)

type _ = base.NoReturn

{service_codes}
"""
#----------------------------

A_SERVICE_CODE ="""type I{service_name} interface {{
{rpc_interfaces}
}}

type {service_name} struct {{
	Sender easygo.IMessageSender
}}

func (self *{service_name}) Init(sender easygo.IMessageSender) {{
	self.Sender = sender
}}

//-------------------------------

{rpc_methods}"""
#----------------------------

A_INTERFACE_ITEM_1 = """\t{method_name}(reqMsg *{request_msg}) *{response_msg}
\t{method_name}_(reqMsg *{request_msg}) (*{response_msg}, easygo.IRpcInterrupt)"""

# rpc 不需 response 的。所以函数也不需要返回值
A_INTERFACE_ITEM_2 = """\t{method_name}(reqMsg *{request_msg})"""

# ------------------------------

A_RPC_METHOD_1 = """func (self *{service_name}) {method_name}(reqMsg *{request_msg}) *{response_msg} {{
	msg, e := self.Sender.CallRpcMethod("{method_name}", reqMsg)
	easygo.PanicError(e)
	if msg == nil {{
		return nil
	}}
	return msg.(*{response_msg})
}}

func (self *{service_name}) {method_name}_(reqMsg *{request_msg}) (*{response_msg}, easygo.IRpcInterrupt) {{
	msg, e := self.Sender.CallRpcMethod("{method_name}", reqMsg)
	if msg == nil {{
		return nil, e
	}}
	return msg.(*{response_msg}), e
}}"""

# rpc 不需 response 的。所以函数也不需要返回值
A_RPC_METHOD_2 = """func (self *{service_name}) {method_name}(reqMsg *{request_msg}) {{
	self.Sender.CallRpcMethod("{method_name}", reqMsg)
}}"""

if __name__ == "__main__":
	args = unpack_args(sys.argv[1:])
	src_path = args["src"] # 源目录
	dest_path = args["dest"] # 目标目录
	rpc_map_file = args["rpc_map_file"] #
	
	up_rpc_items,down_rpc_items = [], []
	for src_name in os.listdir(src_path): # 不再使用 os.walk() 遍历，没有递归的需求
		if not src_name.endswith(EXT_SRC):
			continue
		name, _ = src_name.split(".")
		srcFile = os.path.join(src_path, src_name)
		destFile = dest_path + name + "_rpc" + EXT_DEST

		service_names, ups ,downs = deal_a_proto_file(srcFile, destFile)
		up_rpc_items.extend(ups)
		down_rpc_items.extend(downs)

	up_rpc_items = '\n'.join(up_rpc_items)
	down_rpc_items = '\n'.join(down_rpc_items)
	schema_map = SCHEMA_MAP.format(rpc_map_file=rpc_map_file,up_rpc_items=up_rpc_items,down_rpc_items=down_rpc_items)
	file_path = dest_path + rpc_map_file + EXT_DEST
	
	f = open(file_path, "w",encoding="utf-8")
	f.write(schema_map)
	f.close()
