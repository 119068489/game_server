@echo off
rem ���ļ������� utf-8 ���룬��Ȼ������ʾ�����룬��Ϊ���� Windows ������

set src_path=.\topic_hall\
set dest_path=..\pb\topic_hall\
set rpc_map_file=topic_hall
@echo on

python make_rpc_service.py src=%src_path% dest=%dest_path% rpc_map_file=%rpc_map_file%

@echo off
rem ���� *_rpc.go �ɹ�
rem ==================================================
@echo on

call protoc.exe -I=..\..\github.com\akqp2019\protobuf\protobuf\ -I=..\..\ -I=.\topic_hall -I=..\easygo\base\ -I=.\share_message -I=.\ --gogofast_out=..\pb\topic_hall topic_hall\*.proto


set src_path=..\pb\topic_hall\

python deal_pb_import.py src=%src_path% 
python deal_pb_import.py src=%src_path%
python ..\format_code.py %dest_path%
@pause

@goto end
	---------------------------------------------------------------
	���ȫ��*.proto�ļ�(������import��*.proto�ļ�)���ڵ�ǰĿ¼��ʱ
	����Ҫ-I=��ָʾproto�ļ�����
	�� protoc.exe --python_out=../pb2  *.proto
	---------------------------------------------------------------
	�����import��proto�ļ����ڵ�ʱĿ¼ʱ,��Ҫ��-I=��ָʾ"ͷ�ļ�"����
	ͬʱҪ��ǰĿ¼�µ�protoҲ��Ҫ��-I=��ָ��Ŀ¼
	�������Ҫ����-I=�ֱ�ָ��ͷ�ļ�����proto�ļ�
	�� protoc.exe  -I=../rpc -I=./ --python_out=../pb2  *.proto
	---------------------------------------------------------------
	�����*.proto��д��import "abu/rpc/void.proto";ʱ
	ִ��protoc.exe����*_pb2.pyʱ,����abu/rpc/void.proto�Ƿ����
	���"abu/rpc/void.proto"�Ƿ�����Ǵ�-I=ָ����·���²���
	���Ŀ��*_pb2.py�ļ���Ҳ������import abu.rpc.void_pb2���
	��*.proto�Ͽ���ֱ��ʹ��void.proto�����msg,��������void.xxx������ǰ׺
	---------------------------------------------------------------
:end