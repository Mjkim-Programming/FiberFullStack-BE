import { useState, useEffect } from "react"

const API_URL = "http://localhost:4000/user";

type User = {
  id: number,
  name: string,
  age: number,
};

function App() {
  const [users, setUsers] = useState<User[]>([]);
  const [name, setName] = useState("");
  const [age, setAge] = useState("");

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const res = await fetch(API_URL);
      const data = await res.json();
      setUsers(data);
    } catch (err) {
      console.error("Error fetching users : ", err);
    }
  };

  const addUser = async () => {
    if (!name || !age) {
      alert("이름과 나이를 입력하세요!");
      return;
    }

    try {
      const res = await fetch(API_URL, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({name, age: parseInt(age)}),
      });

      if (!res.ok) throw new Error("Failed to add user");

      setName("");
      setAge("");
      fetchUsers();
    } catch (err) {
      console.error("Error adding user : ", err);
    }
  }

  return (
    <div style={{ padding: "20px", maxWidth: "600px", margin: "auto" }}>
      <div style={{ padding: "20px", maxWidth: "600px", margin: "auto" }}>
        <h1>📌 User Management</h1>

        <div>
          <input
            type="text"
            placeholder="이름 입력"
            value={name}
            onChange={(e) => setName(e.target.value)}
            style={{ marginRight: "10px" }}
          />
          <input
            type="number"
            placeholder="나이 입력"
            value={age}
            onChange={(e) => setAge(e.target.value)}
            style={{ marginRight: "10px" }}
          />
          <button onClick={addUser}>추가</button>
        </div>

        <h2>👥 사용자 목록</h2>
        <ul>
          {users.length > 0 ? (
            users.map((user) => (
              <li key={user.id}>
                {user.name} ({user.age}세)
              </li>
            ))
          ) : (
            <p>사용자가 없습니다.</p>
          )}
        </ul>
      </div>
    </div>
  );
}

export default App