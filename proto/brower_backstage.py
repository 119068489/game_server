import os
import platform

def BuildProto():
    """
    build .proto ---> .pb.go
    """
    src = "./brower_backstage_proto/"
    dest = "../pb/brower_backstage/"
    rpc_map_file="brower_backstage"
    cmd = "python3 make_rpc_service.py src=%s dest=%s rpc_map_file=%s" % (src, dest, rpc_map_file)
    res1 = os.popen(cmd).read()
    print(res1)
    system = platform.system().lower()
    if system == "windows":
        cmd2 = "protoc.exe -I=../../github.com/akqp2019/protobuf/protobuf/ -I=../../ -I=./brower_backstage_proto -I=../easygo/base/ -I=./share_message/ --gogofast_out=../pb/brower_backstage/ brower_backstage_proto/*.proto"
    else:
        cmd2 = "protoc -I=../../ -I=./share_message/ -I=./brower_backstage_proto -I=../easygo/base/ --gogofast_out=../pb/brower_backstage/ brower_backstage_proto/*.proto"
    res2 = os.popen(cmd2).read()
    print(res2)
    src_path="../pb/brower_backstage/"
    cmd3 = "python3 deal_pb_import.py src=%s" % src_path
    res3 = os.popen(cmd3).read()
    print(res3)

if __name__ == "__main__":
    BuildProto()