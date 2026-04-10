from smolagents import CodeAgent, OpenAIServerModel

model = OpenAIServerModel(model_id="{{MODEL}}")

agent = CodeAgent(
    tools=[],
    model=model,
)

if __name__ == "__main__":
    import sys
    message = sys.argv[1] if len(sys.argv) > 1 else ""
    response = agent.run(message)
    print(response)
