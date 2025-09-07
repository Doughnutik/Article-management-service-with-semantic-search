import numpy as np
import psycopg


class AsyncPostgresEmbeddings:
    def __init__(self, dsn: str):
        """
        dsn: строка подключения к Postgres
        """
        self.dsn = dsn
        self.conn: psycopg.AsyncConnection | None = None

    async def connect(self):
        """Создание асинхронного подключения"""
        self.conn = await psycopg.AsyncConnection.connect(self.dsn)
        await self.conn.set_autocommit(True)

    async def insert_embedding(self, id: int, embedding: np.ndarray):
        """Вставка новой пары id, embedding"""
        async with self.conn.cursor() as cur:
            await cur.execute(
                "INSERT INTO embeddings (id, embedding) VALUES (%s, %s)",
                (id, embedding)
            )

    async def delete_embedding(self, id: int):
        """
        Удаление embedding по id.
        Возвращает количество удалённых строк (0, если такого id не было).
        """
        async with self.conn.cursor() as cur:
            await cur.execute("DELETE FROM embeddings WHERE id = %s", (id,))

    async def update_embedding(self, id: int, embedding: np.ndarray):
        """
        Обновление embedding для существующего id.
        Возвращает количество обновлённых строк (0, если такого id не было).
        """
        async with self.conn.cursor() as cur:
            await cur.execute(
                "UPDATE embeddings SET embedding = %s WHERE id = %s",
                (embedding, id)
            )

    async def search_similar(self, query_vector: np.ndarray, k: int = 3) -> list[int]:
        """
        Поиск k ближайших embedding по косинусной близости.
        Возвращает список id.
        """
        async with self.conn.cursor() as cur:
            await cur.execute(
                """
                SELECT id
                FROM embeddings
                ORDER BY embedding <=> %s
                LIMIT %s
                """,
                (query_vector, k)
            )
            rows = await cur.fetchall()
        return [r[0] for r in rows]