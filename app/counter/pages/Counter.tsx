import { signal } from "@preact/signals"

const count = signal(0)

export default function Counter() {
    return (
        <main className="container">
            <h1>Signal Counter</h1>
            <p>This state is global (outside the component).</p>
            <h2>Count: {count.value}</h2>
            <div style={{ display: "flex", gap: "1rem" }}>
                <button onClick={() => count.value++}>+</button>
                <button onClick={() => count.value--}>-</button>
            </div>
        </main>
    )
}
