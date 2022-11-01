from os.path import dirname, abspath
import sys

#3つ上のディレクトリの絶対パスを取得し、sys.pathに登録する
parent_dir = dirname(dirname(dirname(abspath(__file__)))) # 追加
if parent_dir not in sys.path: # 追加
    sys.path.append(parent_dir) # 追加



import grpc
import logging
from concurrent import futures
from api.v1 import inference_pb2
from api.v1 import inference_pb2_grpc
from inference import model_output

#from api.v1 import inference_pb2 as api_dot_v1_dot_inference__pb2

class InferenceServicer(inference_pb2_grpc.InferenceServicer):

    def InferenceCall(self, request, context):
        predicted = model_output(request.text)

        return inference_pb2.Output_data_for_inference(text=request.text, label=predicted)


def serve():
  server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
  inference_pb2_grpc.add_InferenceServicer_to_server(
      InferenceServicer(), server)
  server.add_insecure_port('127.0.0.1:8403')
  server.start()
  server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig()
    serve()