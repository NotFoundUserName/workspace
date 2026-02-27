(() => {
  const chat = document.getElementById("chat");
  const composer = document.getElementById("composer");
  const messageInput = document.getElementById("message");
  const sendBtn = document.getElementById("send");

  function nowTime() {
    try {
      return new Date().toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    } catch {
      return "";
    }
  }

  function append(role, text) {
    const row = document.createElement("div");
    row.className = `row ${role}`;

    const bubble = document.createElement("div");
    bubble.className = "bubble";
    bubble.textContent = text;

    const meta = document.createElement("div");
    meta.className = "meta";
    meta.textContent = `${role === "user" ? "你" : "Bot"} · ${nowTime()}`;

    const wrap = document.createElement("div");
    wrap.appendChild(bubble);
    wrap.appendChild(meta);

    row.appendChild(wrap);
    chat.appendChild(row);
    chat.scrollTop = chat.scrollHeight;
  }

  async function sendMessage(text) {
    sendBtn.disabled = true;
    messageInput.disabled = true;
    messageInput.classList.remove("error");

    try {
      const res = await fetch("/api/chat", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ message: text }),
      });

      if (!res.ok) {
        const errText = await res.text().catch(() => "");
        throw new Error(errText || `HTTP ${res.status}`);
      }

      const data = await res.json();
      append("bot", (data && data.reply) || "（无响应）");
    } catch (e) {
      append("bot", `请求失败：${e && e.message ? e.message : "unknown error"}`);
    } finally {
      sendBtn.disabled = false;
      messageInput.disabled = false;
      messageInput.focus();
    }
  }

  (async () => {
    try {
      const res = await fetch("/api/info", { method: "GET" });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const info = await res.json();
      const greet = (info && info.greeting) || "你好！今天我能帮你什么忙？";
      append("bot", greet);
    } catch {
      append("bot", "你好！今天我能帮你什么忙？");
    } finally {
      messageInput.focus();
    }
  })();

  composer.addEventListener("submit", async (ev) => {
    ev.preventDefault();
    const text = (messageInput.value || "").trim();
    if (!text) {
      messageInput.classList.add("error");
      messageInput.focus();
      return;
    }
    messageInput.value = "";
    append("user", text);
    await sendMessage(text);
  });
})();

