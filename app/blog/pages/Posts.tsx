import { useEffect, useState } from "uca/ui"
import { Post } from "uca/srv"

export default function Posts() {
    const [title, setTitle] = useState("")
    const [content, setContent] = useState("")

    useEffect(() => {
        Post.GET.fetch()
    }, [])

    const submitPost = (e: any) => {
        e.preventDefault()
        Post.POST({ title, content }).then(() => {
            setTitle("")
            setContent("")
        })
    }

    return (
        <main className="container">
            <h1>My Blog</h1>
            <form onSubmit={submitPost}>
                <input placeholder="Title" value={title} onInput={e => setTitle((e.target as HTMLInputElement).value)} required />
                <textarea placeholder="Content" value={content} onInput={e => setContent((e.target as HTMLTextAreaElement).value)} required />
                <button type="submit">Publish</button>
            </form>

            <hr />

            {Post.GET.data.value.map(p => (
                <article key={p.id}>
                    <h3>{p.title}</h3>
                    <p>{p.content}</p>
                </article>
            ))}
        </main>
    )
}
