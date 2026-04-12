import { useEffect, useRef, useState } from "preact/hooks"

export default function Chat() {
    const [messages, setMessages] = useState([])
    const [input, setInput] = useState("")
    const bottomRef = useRef(null)

    useEffect(() => {
        fetch("/api/History")
            .then(res => res.json())
            .then(data => {
                if (Array.isArray(data)) {
                    setMessages(data)
                } else {
                    setMessages([])
                }
            })
    }, [])

    useEffect(() => {
        bottomRef.current?.scrollIntoView({ behavior: "smooth" })
    }, [messages])

    const saveMessage = (sender, content) => {
        fetch("/api/History", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ sender, content })
        })
    }

    const sendMessage = () => {
        if (!input.trim()) return

        const userMsg = { sender: "user", content: input, timestamp: new Date().toISOString() }
        setMessages(prev => [...prev, userMsg])
        saveMessage("user", input)
        setInput("")

        const agentMsg = { sender: "agent", content: "", timestamp: new Date().toISOString() }
        setMessages(prev => [...prev, agentMsg])

        const es = new EventSource("/api/chat/Assistant?message=" + encodeURIComponent(userMsg.content))

        es.onmessage = (e) => {
            if (e.data === "[DONE]") {
                es.close()
                return
            }
            agentMsg.content = e.data
            setMessages(prev => [...prev.slice(0, -1), { ...agentMsg }])
            saveMessage("agent", e.data)
        }

        es.onerror = () => es.close()
    }

    return (
        <div style={{
            display: "flex",
            flexDirection: "column",
            height: "100vh",
            maxWidth: "98%",
            margin: "0 auto",
            padding: "1rem"
        }}>
            <div style={{
                flex: 1,
                overflowY: "auto",
                display: "flex",
                flexDirection: "column",
                gap: "1.5rem",
                paddingBottom: "1rem"
            }}>
                {messages.map((msg, i) => (
                    <div key={i} style={{ display: "flex", flexDirection: "column", gap: "0.25rem" }}>
                        <small style={{
                            fontWeight: "600",
                            color: msg.sender === "user" ? "var(--pico-primary)" : "var(--pico-muted-color)",
                            textTransform: "capitalize"
                        }}>
                            {msg.sender}
                        </small>
                        <p style={{ margin: 0, lineHeight: "1.6" }}>{msg.content}</p>
                    </div>
                ))}
                <div ref={bottomRef} />
            </div>
            <div style={{
                display: "flex",
                gap: "0.5rem",
                paddingTop: "1rem",
                borderTop: "1px solid var(--pico-muted-border-color)"
            }}>
                <input
                    style={{ flex: 1, margin: 0 }}
                    value={input}
                    onInput={e => setInput((e.target as HTMLInputElement).value)}
                    onKeyDown={e => e.key === "Enter" && sendMessage()}
                    placeholder="Type a message..."
                />
                <button onClick={sendMessage} style={{ margin: 0 }}>Send</button>
            </div>
        </div>
    )
}
