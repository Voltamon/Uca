import { useState } from "uca/ui"

export default function Calc() {
    const [expr, setExpr] = useState("")
    const [res, setRes] = useState("")

    const calculate = () => {
        try {
            setRes(eval(expr).toString())
        } catch {
            setRes("Error")
        }
    }

    return (
        <main className="container">
            <h1>Calculator</h1>
            <input 
                value={expr} 
                onInput={e => setExpr((e.target as HTMLInputElement).value)} 
                placeholder="2 + 2"
            />
            <button onClick={calculate}>=</button>
            {res && <h2>Result: {res}</h2>}
        </main>
    )
}
