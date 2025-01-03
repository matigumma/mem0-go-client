import os
from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel, Field
from typing import List, Optional, Dict, Any
from mem0client import Mem0Client, Message, StoreOptions, GetMemoriesOptions, SearchMemoriesOptions

# Load API key from environment variable
API_KEY = os.getenv('MEM0_API_KEY')
if not API_KEY:
    raise ValueError("MEM0_API_KEY environment variable is required")

# Initialize Mem0 client
mem0_client = Mem0Client(API_KEY)

app = FastAPI(
    title="Mem0 Memory Management API",
    description="API wrapper for Mem0 memory management service",
    version="0.1.0"
)

class MessageModel(BaseModel):
    role: str
    content: str

class StoreMemoryRequest(BaseModel):
    messages: List[MessageModel]
    user_id: Optional[str] = None
    agent_id: Optional[str] = None
    run_id: Optional[str] = None
    metadata: Optional[Dict[str, Any]] = None

class SearchMemoriesRequest(BaseModel):
    query: str
    user_id: Optional[str] = None
    agent_id: Optional[str] = None
    limit: Optional[int] = 10

@app.post("/memories/store")
async def store_memory(request: StoreMemoryRequest):
    """
    Store a new memory with optional metadata and identifiers
    """
    try:
        # Convert Pydantic messages to Mem0 messages
        messages = [
            Message(role=msg.role, content=msg.content) 
            for msg in request.messages
        ]
        
        store_options = StoreOptions(
            messages=messages,
            user_id=request.user_id,
            agent_id=request.agent_id,
            run_id=request.run_id,
            metadata=request.metadata or {}
        )
        
        result = mem0_client.Store(store_options)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/memories")
async def get_memories(
    user_id: Optional[str] = Query(None),
    agent_id: Optional[str] = Query(None),
    app_id: Optional[str] = Query(None),
    run_id: Optional[str] = Query(None)
):
    """
    Retrieve memories with optional filtering
    """
    try:
        options = GetMemoriesOptions(
            user_id=user_id,
            agent_id=agent_id,
            app_id=app_id,
            run_id=run_id
        )
        
        memories = mem0_client.GetMemories(options)
        return memories
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/memories/search")
async def search_memories(request: SearchMemoriesRequest):
    """
    Perform semantic search on memories
    """
    try:
        search_options = SearchMemoriesOptions(
            query=request.query,
            user_id=request.user_id,
            agent_id=request.agent_id,
            limit=request.limit
        )
        
        results = mem0_client.SearchMemories(search_options)
        return results
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)