@echo off
rem ���ļ������� utf-8 ���룬��Ȼ������ʾ�����룬��Ϊ���� Windows ������


rem ==================================================
@echo on
call protoc.exe -I=..\..\ -I=.\share_message --gogofast_out=..\pb\share_message share_message\*.proto 

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
