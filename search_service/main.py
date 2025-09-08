import asyncio
import grpc
from grpc_python import grpc_pb2, grpc_pb2_grpc
from embeddings.emb import ArticleEmbedder
from storage.postgres_embeddings import AsyncPostgresEmbeddings


class SearchService(grpc_pb2_grpc.SearchServiceServicer):
    def __init__(self, db_dsn: str, model_name: str = "jinaai/jina-embeddings-v3", dim: int = 1024):
        self.model = ArticleEmbedder(model_name=model_name, dim=dim)
        self.db = AsyncPostgresEmbeddings(db_dsn)

    async def IndexArticle(self, request, context):
        text = f"{request.title}. {request.content}"
        embedding = self.model.embed_article(text)
        await self.db.insert_embedding(request.id, embedding)
        print(f"Indexed article {request.id}")
        return grpc_pb2.ArticleEmbeddingResponse()

    async def SemanticSearch(self, request, context):
        query_emb = self.model.embed_query(request.query)
        nearest_ids = await self.db.search_similar(query_emb, k=request.limit or 3)
        print(f"Search query: {request.query}, found ids: {nearest_ids}")
        return grpc_pb2.SearchResponse(article_ids=nearest_ids)


async def serve():
    username = "artem"
    password = "1234"
    db_dsn = f"postgresql://{username}:{password}@localhost:5432/embeddings"
    service = SearchService(db_dsn)
    await service.db.connect()

    server = grpc.aio.server()
    grpc_pb2_grpc.add_SearchServiceServicer_to_server(service, server)
    server.add_insecure_port('[::]:50051')
    await server.start()
    print("Async gRPC server started on port 50051")
    await server.wait_for_termination()


if __name__ == '__main__':
    asyncio.run(serve())