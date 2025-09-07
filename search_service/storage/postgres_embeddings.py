import psycopg2
import numpy as np

class PostgresEmbeddings:
    def __init__(self, dsn: str):
        """
        dsn: строка подключения к Postgres
        """
        self.conn = psycopg2.connect(dsn)
        self.conn.autocommit = True

    def insert_embedding(self, id: int, embedding: np.ndarray):
        """Вставка новой пары id, embedding"""
        with self.conn.cursor() as cur:
            cur.execute(
                "INSERT INTO embeddings (id, embedding) VALUES (%s, %s)",
                (id, embedding.tolist())
            )

    def delete_embedding(self, id: int):
        """Удаление embedding по id"""
        with self.conn.cursor() as cur:
            cur.execute("DELETE FROM embeddings WHERE id = %s", (id,))

    def update_embedding(self, id: int, embedding: np.ndarray):
        """Обновление embedding для существующего id"""
        with self.conn.cursor() as cur:
            cur.execute(
                "UPDATE embeddings SET embedding = %s WHERE id = %s",
                (embedding.tolist(), id)
            )

    def search_similar(self, query_vector: np.ndarray, k: int = 3) -> list[int]:
        """
        Поиск k ближайших embedding по косинусной близости
        Возвращает список id
        """
        with self.conn.cursor() as cur:
            cur.execute(
                """
                SELECT id
                FROM embeddings
                ORDER BY embedding <=> %s
                LIMIT %s
                """,
                (query_vector.tolist(), k)
            )
            rows = cur.fetchall()
        return [r[0] for r in rows]