# Security and Isolation

- Central AI server with **stateless requests per user**
- Memory / RAG stored in backend databases scoped by `user_id` or `org_id`
- No cross-user data leakage
- TLS for all network communications
- Optional enterprise namespaces for strict isolation
- Logging includes user context for debugging without exposing others’ data
