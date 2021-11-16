import os
import platform


def BuildProto():
    """
    build .proto ---> .pb.go
    """
    src = "./h5_server_proto/h5_wish/"
    dest = "../pb/h5_wish/"
    rpc_map_file = "h5_wish"
    cmd = "python3 make_rpc_service.py src=%s dest=%s rpc_map_file=%s" % (src, dest, rpc_map_file)
    res1 = os.popen(cmd).read()
    print(res1)
    system = platform.system().lower()
    if system == "windows":
        cmd2 = "protoc.exe -I=../../github.com/akqp2019/protobuf/protobuf/ -I=../../ -I=./h5_server_proto/h5_wish -I=../easygo/base/ -I=./share_message --gogofast_out=../pb/h5_wish h5_server_proto/h5_wish/*.proto"
    else:
        cmd2 = "protoc -I=../../github.com/akqp2019/protobuf/protobuf/ -I=../../ -I=./h5_server_proto/h5_wish -I=../easygo/base/ -I=./share_message --gogofast_out=../pb/h5_wish h5_server_proto/h5_wish/*.proto"
    res2 = os.popen(cmd2).read()
    print(res2)
    src2 = "../pb/h5_wish/"
    cmd3 = "python3 deal_pb_import.py src=%s" % src2
    res3 = os.popen(cmd3).read()
    print(res3)


if __name__ == "__main__":
    BuildProto()