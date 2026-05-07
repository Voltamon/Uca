import { useEffect, useRef, useState } from "uca/ui"
import { History } from "uca/srv"
import { TaskBot } from "uca/ai"

export default function Chat() {
    const [input, setInput] = useState("")
    const bottomRef = useRef(null)

    useEffect(() => {
        History.GET.fetch()
    }, [])

    useEffect(() => {
        bottomRef.current?.scrollIntoView({ behavior: "smooth" })
    }, [History.GET.data.value])

    const sendMessage = () => {
        if (!input.trim()) return
        const text = input
        setInput("")

        // Persist user message - this will auto-refresh History.GET.data
        History.POST({ sender: "user", content: text })

        // Send to AI
        TaskBot.chat(text, (data) => {
            // After AI is done, we might want to persist it too
            // Note: In a real app, you'd save the AI response.
            // For now, let's just save it.
            if (data === "[DONE]") {
                // Actually, the chat helper doesn't give us the final string easily here
            }
        })
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
                {History.GET.data.value.map((msg, i) => (
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
