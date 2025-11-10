# Tech Stack

## Backend / CLI (MVP)

- **Language**: Go  
- **Libraries/Tools**:
  - `cobra` → CLI framework
  - `playwright-go` → Browser automation
  - `viper` → Configuration management
  - `zerolog` → Structured logging
- **AI Integration**: Hugging Face Mistral-7B via inference endpoint  
- **Execution**: Local machine (stateless AI requests per user)  
- **Config**: `config.yaml` for optional selectors  

## DevOps / AI Infrastructure

- **Containerization**: Docker for portability  
- **Cloud Providers**: AWS / GCP for SaaS  
- **Optional Tools**: Terraform / Ansible for enterprise deployment  
- **Vector DB**: PostgreSQL + pgvector / Pinecone / Weaviate (RAG)  

## Frontend / SaaS

- **Frameworks**: React + Next.js + TypeScript + Tailwind CSS  
- **Data Fetching**: TanStack Query  
- **Features**: Authentication, multi-user support, dashboards, test result visualization  
- **Integration**: Optionally connect with Figma for UI comparison  
