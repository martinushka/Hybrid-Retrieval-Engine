# Hybrid Retrieval Engine

A **Go-based backend** for intelligent product search, combining lexical search, semantic retrieval, hybrid ranking, and RAG-powered answer generation.

The system is designed to retrieve the most relevant products for a user query, rank them using multiple signals, and provide the retrieved context to an LLM for generating grounded answers.

## What the System Does

* normalizes and tokenizes text;
* performs lexical search across product titles, categories, and descriptions;
* calculates relevance scores for retrieved products;
* performs semantic search using embeddings and vector similarity;
* combines lexical and semantic results through hybrid ranking;
* retrieves the Top-K most relevant documents;
* builds context for an LLM;
* returns generated answers together with their sources.

## Architecture

```text
User Query
    │
    ▼
Query Processing
    │
    ├─────────────────┐
    ▼                 ▼
Lexical Search   Semantic Search
    │                 │
    └────────┬────────┘
             ▼
       Hybrid Ranking
             │
             ▼
           Top-K
             │
             ▼
     Context Construction
             │
             ▼
            LLM
             │
             ▼
      Answer + Sources
```

## Core Components

**Text Processing** — text normalization and tokenization for queries and documents.

**Lexical Search** — keyword-based retrieval across product titles, categories, and descriptions using relevance scoring.

**Semantic Search** — retrieval based on semantic similarity using text embeddings and vector search.

**Hybrid Ranking** — combines lexical and semantic relevance signals to improve retrieval quality.

**RAG Pipeline** — constructs a context from retrieved documents and uses an LLM to generate grounded answers.

**Sources** — preserves the documents used to generate each answer, improving transparency and traceability.

## Project Direction

The project evolves from a classical **lexical search engine** into a **hybrid retrieval system** and ultimately a full **RAG pipeline**, focusing on retrieval quality, ranking algorithms, vector search, and production-oriented backend architecture.
