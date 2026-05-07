import { useState } from "uca/ui"
import { Forecaster } from "uca/ai"

export default function Weather() {
    const [city, setCity] = useState("")
    const [report, setReport] = useState("")
    const [loading, setLoading] = useState(false)

    const ask = () => {
        setLoading(true)
        setReport("")
        Forecaster.chat(city, (chunk) => {
            setReport(prev => prev + chunk)
        }).then(() => setLoading(false))
    }

    return (
        <main className="container">
            <h1>Weather AI</h1>
            <input 
                placeholder="Enter city..." 
                value={city} 
                onInput={e => setCity((e.target as HTMLInputElement).value)} 
            />
            <button onClick={ask} disabled={loading}>Get Forecast</button>
            {report && (
                <article>
                    <p>{report}</p>
                </article>
            )}
        </main>
    )
}
