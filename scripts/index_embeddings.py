import os

import psycopg
from sentence_transformers import SentenceTransformer


DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgres://localhost/ios_rag",
)

MODEL_NAME = "intfloat/multilingual-e5-small"

model = SentenceTransformer(MODEL_NAME)

with psycopg.connect(DATABASE_URL) as conn:
    with conn.cursor() as cur:
        cur.execute("""
            SELECT id, title, description, category
            FROM products
            ORDER BY id
        """)

        products = cur.fetchall()

        for product_id, title, description, category in products:
            text = (
                f"passage: {title}. "
                f"Category: {category}. "
                f"{description}"
            )

            embedding = model.encode(
                text,
                normalize_embeddings=True,
            )

            vector = "[" + ",".join(str(float(x)) for x in embedding) + "]"

            cur.execute(
                """
                UPDATE products
                SET embedding = %s::vector
                WHERE id = %s
                """,
                (vector, product_id),
            )

            print(f"indexed product {product_id}: {title}")

    conn.commit()

print(f"indexed {len(products)} products")
