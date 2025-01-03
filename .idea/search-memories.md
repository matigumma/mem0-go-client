# Mem0 Search Memories API Documentation

## Overview

The Mem0 Search Memories API allows you to perform semantic searches on stored memories. This API is useful for retrieving specific memories based on a query and various filtering options.

## Endpoint

**POST** `/v1/memories/search/`

## Authorization

- **Type**: API key
- **Header**: `Authorization`
- **Format**: `Token your_api_key`

## Request Body

The request body should be in JSON format and include the following fields:

- **query** (string, required): The search query to find in the memory. Minimum length is 1.
- **agent_id** (string | null): The agent ID associated with the memory.
- **user_id** (string | null): The user ID associated with the memory.
- **app_id** (string | null): The app ID associated with the memory.
- **run_id** (string | null): The run ID associated with the memory.
- **metadata** (object | null): Additional metadata associated with the memory.
- **top_k** (integer, default: 10): The number of top results to return.
- **fields** (array of strings): A list of field names to include in the response. If not provided, all fields will be returned.
- **rerank** (boolean, default: false): Whether to rerank the memories.
- **output_format** (string | null): The search method supports two output formats: `v1.0` (default) and `v1.1`.
- **org_id** (string | null): The unique identifier of the organization associated with the memory.
- **project_id** (string | null): The unique identifier of the project associated with the memory.
- **filter_memories** (boolean, default: false): Whether to properly filter the memories according to the input.
- **categories** (array of strings): A list of categories to filter the memories by.
- **only_metadata_based_search** (boolean, default: false): Whether to only search for memories based on metadata.

## Response

The response will be in JSON format and include the following fields:

- **id** (string): Unique identifier for the memory.
- **memory** (string): The content of the memory.
- **input** (array of objects): The conversation input that was used to generate this memory.
  - **role** (enum<string>): The role of the speaker in the conversation. Options: `user`, `assistant`.
  - **content** (string): The content of the message.
- **user_id** (string): The identifier of the user associated with this memory.
- **hash** (string): A hash of the memory content.
- **metadata** (object | null): Additional metadata associated with the memory.
- **created_at** (string): The timestamp when the memory was created.
- **updated_at** (string): The timestamp when the memory was last updated.

## Status Codes

- **200**: Successful retrieval of memories.
- **400**: Bad request, typically due to missing or invalid parameters.
