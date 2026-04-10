import { useEffect, useRef, useState } from "preact/hooks";

export default function Chat() {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const bottomRef = useRef(null);

  useEffect(() => {
    fetch("/api/History")
      .then((res) => res.json())
      .then((data) => setMessages(data));
  }, []);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const sendMessage = () => {
    if (!input.trim()) return;

    const userMsg = {
      sender: "user",
      content: input,
      timestamp: new Date().toISOString(),
    };
    setMessages((prev) => [...prev, userMsg]);
    setInput("");

    const es = new EventSource(
      "/api/chat/Assistant?message=" + encodeURIComponent(input),
    );
    let agentMsg = {
      sender: "agent",
      content: "",
      timestamp: new Date().toISOString(),
    };
    setMessages((prev) => [...prev, agentMsg]);

    es.onmessage = (e) => {
      agentMsg.content += e.data;
      setMessages((prev) => [...prev.slice(0, -1), { ...agentMsg }]);
    };

    es.onerror = () => es.close();
  };

  return (
    <div class="chat">
      <div class="messages">
        {messages.map((msg, i) => (
          <div key={i} class={"message " + msg.sender}>
            <span>{msg.content}</span>
          </div>
        ))}
        <div ref={bottomRef} />
      </div>
      <div class="input-area">
        <input
          value={input}
          onInput={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && sendMessage()}
          placeholder="Type a message..."
        />
        <button onClick={sendMessage}>Send</button>
      </div>
    </div>
  );
}
