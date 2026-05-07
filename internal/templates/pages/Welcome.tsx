import { useEffect, useRef } from "uca/ui"
import { User } from "uca/srv"

export default function Welcome() {
    const inputRef = useRef(null)

    useEffect(() => {
        User.GET.fetch().then(() => {
            if (User.GET.data.value.length > 0) {
                window.location.href = "/chat"
            }
        })
    }, [])

    const handleSubmit = (e: any) => {
        e.preventDefault()
        const name = inputRef.current.value
        if (!name.trim()) return

        User.POST({ name, role: "user" })
            .then(() => {
                window.location.href = "/chat"
            })
    }

    return (
        <main className="container" style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "100vh" }}>
            <article style={{ width: "100%", maxWidth: "400px" }}>
                <header>
                    <h1 style={{ textAlign: "center", margin: 0 }}>Uca</h1>
                    <p style={{ textAlign: "center", opacity: 0.6 }}>Build anything. Fast.</p>
                </header>
                <form onSubmit={handleSubmit}>
                    <label>What's your name?</label>
                    <input 
                        ref={inputRef} 
                        placeholder="e.g. Satoshi" 
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
