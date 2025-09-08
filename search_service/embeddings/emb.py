import numpy as np
from sentence_transformers import SentenceTransformer

class ArticleEmbedder:
    """
    Генерация эмбеддингов для научных статей и запросов.
    """

    def __init__(self, model_name: str = "jinaai/jina-embeddings-v3", dim: int = 1024):
        self.model = SentenceTransformer(model_name, trust_remote_code=True)
        self.dim = dim

    def embed_article(self, text: str) -> np.ndarray:
        """
        Генерация эмбеддинга для статьи
        """
        return self.model.encode(text, task="retrieval.passage", normalize_embeddings=True, convert_to_numpy=True, truncate_dim=self.dim)

    def embed_query(self, text: str) -> np.ndarray:
        """
        Генерация эмбеддинга для запроса
        """
        return self.model.encode(text, task="retrieval.query", normalize_embeddings=True, convert_to_numpy=True, truncate_dim=self.dim)
    
    def cosine_similarity(self, first_emb: np.ndarray, second_emb: np.ndarray):
        scal_prod = first_emb @ second_emb
        v1 = first_emb @ first_emb
        v2 = second_emb @ second_emb
        norm = np.sqrt(v1 * v2)
        return scal_prod / norm


if __name__ == "__main__":
    embedder = ArticleEmbedder()

    articles = [
        "Научная статья о квантовой физике и её экспериментах",
        "Статья о машинном обучении и нейросетях",
        "Обзор современных технологий в биоинформатике",
        "Квантовая механика",
        "Что такое квантовый мир и какая в нём механика жизни",
        "what is the quantum mechanics",
        "КаК РешаТь УравНения"
    ]
    query = "How to solve equations"

    article_embeddings = [embedder.embed_article(a) for a in articles]
    query_embedding = embedder.embed_query(query)
    
    for i in article_embeddings:
        print(embedder.cosine_similarity(i, query_embedding))