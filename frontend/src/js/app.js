const BACKEND_URL = "BACKEND_URL"; 

fetch(`${BACKEND_URL}/api/messages`)
  .then(response => response.json())
  .then(data => {
    const list = document.getElementById("message-list");
    data.forEach(msg => {
      const li = document.createElement("li");
      li.textContent = msg.message;
      list.appendChild(li);
    });
  })
  .catch(err => {
    console.error("Failed to fetch messages:", err);
    const list = document.getElementById("message-list");
    const li = document.createElement("li");
    li.textContent = "Error fetching data from backend.";
    list.appendChild(li);
  });