import { useState, useEffect } from "uca/ui"
import { Task } from "uca/srv"

export default function Home() {
    const [newTitle, setNewTitle] = useState("")

    useEffect(() => {
        Task.GET.fetch()
    }, [])

    const addTask = (e: any) => {
        e.preventDefault()
        if (!newTitle.trim()) return
        Task.POST({ title: newTitle, done: false }).then(() => {
            setNewTitle("")
        })
    }

    return (
        <main className="container">
            <h2>Todo List</h2>
            <form onSubmit={addTask}>
                <input 
                    value={newTitle} 
                    onInput={e => setNewTitle((e.target as HTMLInputElement).value)}
                    placeholder="New Task"
                />
                <button type="submit">Add</button>
            </form>

            {Task.GET.loading.value ? <p>Loading...</p> : (
                <ul>
                    {Task.GET.data.value.map(t => (
                        <li key={t.id}>
                            <input 
                                type="checkbox" 
                                checked={t.done} 
                                onChange={() => Task.PUT({ id: t.id, done: !t.done })}
                            />
                            {t.title}
                            <button onClick={() => Task.DELETE({ id: t.id })}>x</button>
                        </li>
                    ))}
                </ul>
            )}
        </main>
    )
}
