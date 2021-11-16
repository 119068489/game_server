import os
import platform


def BuildProto():
    """
    build .proto ---> .pb.go
    """
    src = "./client_server_proto/"
    dest = "../pb/client_server/"
    rpc_map_file = "client_server"
    cmd = "python3 make_rpc_service.py src=%s dest=%s rpc_map_file=%s" % (src, dest, rpc_map_file)
    res = os.popen(cmd).read()
    print(res)
    system = platform.system().lower()
    if system == "windows":
        cmd2 = "protoc.exe -I=../../github.com/akqp2019/protobuf/protobuf/ -I=../../ -I=./client_server_proto -I=../easygo/base/ -I=./share_message --gogofast_out=../pb/client_server client_server_proto/*.proto"
    else:
        cmd2 = "protoc -I=../../github.com/akqp2019/protobuf/protobuf/ -I=../../ -I=./client_server_proto -I=../easygo/base/ -I=./share_message --gogofast_out=../pb/client_server client_server_proto/*.proto"
    res2 = os.popen(cmd2).read()
    print(res2)
    src2 = "../pb/client_server"
    cmd3 = "python3 deal_pb_import.py src=%s" % src2
    res3 = os.popen(cmd3).read()
    print(res3)

if __name__ == "__main__":
    BuildProto()
