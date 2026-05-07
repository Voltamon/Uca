from uca.ai import Agent, Message

def get_weather(city: str):
    """Returns the weather for a city."""
    # In a real app, you'd fetch from an API
    if city.lower() == "tokyo": return "Rainy, 18°C"
    return "Sunny, 25°C"

agent = Agent(
    model="github/gpt-4o-mini",
    tools=[get_weather]
)

agent.prompt = f"System: You are a weather forecaster. Use the tool to provide accurate info.\nUser: {Message}"
