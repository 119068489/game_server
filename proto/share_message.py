import os
import platform


def BuildProto():
    """
    build .proto ---> .pb.go
    """
    system = platform.system().lower()
    if system == 'windows':
        excute = 'protoc.exe -I=../../github.com/akqp2019/protobuf/protobuf/ -I=../../ -I=./share_message --gogofast_out=../pb/share_message share_message/*.proto'
    else:
        excute = 'protoc -I=../../github.com/akqp2019/protobuf/protobuf/ -I=../../ -I=./share_message --gogofast_out=../pb/share_message share_message/*.proto'
    result = os.popen(excute).read()
    print(result)

if __name__ == "__main__":
    BuildProto()