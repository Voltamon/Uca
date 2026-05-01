from uca import Agent

assistant = Agent(model_id="{{MODEL}}")

if __name__ == "__main__":
    import sys
    message = sys.argv[1] if len(sys.argv) > 1 else ""
    response = assistant.run(message)
    print(response)
