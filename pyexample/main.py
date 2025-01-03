from mem0 import Memory

config = {
    "version": "v1.1"
}

m = Memory.from_config(config_dict=config)

# For a user
result = m.add("Likes to play cricket on weekends", 
user_id="alice", 
metadata={"category": "hobbies"}, 
)

print("Add Result: ")
print(result)

# messages = [
#    {"role": "user", "content": "Hi, I'm Alex. I like to play cricket on weekends."},
#    {"role": "assistant", "content": "Hello Alex! It's great to know that you enjoy playing cricket on weekends. I'll remember that for future reference."}
# ]
# client.add(messages, user_id="alice")


# Get all memories

all_memories = m.get_all(user_id="alice")
print("Get All memories: ")
print(all_memories)

# Get a single memory by ID
# specific_memory = m.get("bf4d4092-cf91-4181-bfeb-b6fa2ed3061b")


# related_memories = m.search(query="What are Alice's hobbies?", user_id="alice")


result = m.update(memory_id="f5e54022-0a39-49b2-aed1-8c3fc1c57599", data="Likes to play tennis on weekends")

print("Update Result: ")
print(result)

# history = m.history(memory_id="bf4d4092-cf91-4181-bfeb-b6fa2ed3061b")


# # Delete a memory by id
# m.delete(memory_id="bf4d4092-cf91-4181-bfeb-b6fa2ed3061b")
# # Delete all memories for a user
# m.delete_all(user_id="alice")



# m.reset() # Reset all memories
