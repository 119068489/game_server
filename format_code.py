# -*- coding: utf-8 -*-

import sys
import os

POSTFIX = ".go" # 目标文件扩展名

ignore = ('.git', 'proto') # 要忽略的目录



def entry(directory):
	
	for dir_path, dir_names, file_names in os.walk(directory):

		dirs = dir_path.split(os.path.sep)
		if set(dirs) & set(ignore):
			 continue

		for file_name in file_names:
			if not file_name.endswith(POSTFIX):
				continue
			file_path = os.path.join(dir_path, file_name)
			# print(file_path)
			
			cmd = "go fmt " + file_path
			# print(cmd)
			os.system(cmd)
			

if __name__ == "__main__":
	if len(sys.argv)>1:
		directory = sys.argv[1]
	else:
		directory = os.getcwd()
	entry(directory)

