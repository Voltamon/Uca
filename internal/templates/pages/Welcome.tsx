import { useEffect, useRef } from "uca/ui"

export default function Welcome() {
    const inputRef = useRef(null)

    useEffect(() => {
        fetch("/api/User")
            .then(res => {
                if (res.status === 404) return null
                return res.json()
            })
            .then(data => {
                if (data && data.name) {
                    window.location.href = "/chat"
                }
            })
            .catch(() => {})
    }, [])

    const handleSubmit = (e) => {
        e.preventDefault()
        const name = inputRef.current.value
        if (!name.trim()) return

        fetch("/api/User", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name })
        })
        .then(res => res.json())
        .then(() => {
            window.location.href = "/chat"
        })
    }

    return (
        <main style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            minHeight: "100vh"
        }}>
            <article style={{ width: "100%", maxWidth: "400px" }}>
                <header>
                  <h2>Welcome to {{APP_NAME}}</h2>
                    <p>Tell us your name to get started</p>
                </header>
                <form onSubmit={handleSubmit}>
                    <input
                        ref={inputRef}
                        type="text"
                        placeholder="Your name"
                        required
                    />
                    <button type="submit" style={{ width: "100%" }}>
                        Get Started
                    </button>
                </form>
            </article>
        </main>
    )
}
