import { useEffect } from "preact/hooks";

export default function Welcome() {
  useEffect(() => {
    fetch("/api/User")
      .then((res) => res.json())
      .then((data) => {
        if (data && data.name) {
          window.location.href = "/chat";
        }
      });
  }, []);

  const handleSubmit = (e) => {
    e.preventDefault();
    const name = e.target.name.value;
    fetch("/api/User", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name }),
    }).then(() => {
      window.location.href = "/chat";
    });
  };

  return (
    <div class="welcome">
      <h1>Welcome to {{ APP_NAME }}</h1>
      <form onSubmit={handleSubmit}>
        <input
          name="name"
          type="text"
          placeholder="What is your name?"
          required
        />
        <button type="submit">Get Started</button>
      </form>
    </div>
  );
}
