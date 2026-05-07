from uca.ai import Agent, Message
from uca.srv import Todo

agent = Agent(
    model="{{MODEL}}",
    tools=[Todo.All]
)

# You can customize the prompt using the Message placeholder
agent.prompt = f"System: You are a helpful assistant.\nUser: {Message}"
