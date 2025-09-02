import grpc
from concurrent import futures
from grpc_python import grpc_pb2, grpc_pb2_grpc

class SearchService(grpc_pb2_grpc.SearchServiceServicer):
    def IndexArticle(self, request, context):
        #TODO сделать построение эмбеддинга и сохранение его в базу данных
        
        print(f"Indexing article with id = {request.id}, title = {request.title}, content = {request.content}")
        return grpc_pb2.ArticleEmbeddingResponse()

    def SemanticSearch(self, request, context):
        #TODO сделать семантический поиск по эмбеддингам в базе данных
        
        print(f"Searching for: {request.query}")
        return grpc_pb2.SearchResponse(article_ids=[1, 2, 3])

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    grpc_pb2_grpc.add_SearchServiceServicer_to_server(SearchService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()