import { useState, useEffect } from "react";
import logo from "./assets/images/logo-universal.png";
import "./App.css";
import { EventsOn } from "../wailsjs/runtime/runtime";
import { Register } from "../wailsjs/go/auth/AuthService";
import { useNavigate } from "react-router";

function App() {
  const [name, setName] = useState("");

  const navigate = useNavigate();

  useEffect(() => {
    EventsOn("navigate", (path) => {
      navigate(path);
    });
  }, [navigate]);

  return (
    <div id="App">
      <img src={logo} id="logo" alt="logo" />
      <div id="input" className="input-box">
        <input
          id="name"
          className="input"
          onChange={(e) => setName(e.target.value)}
          autoComplete="off"
          name="input"
          type="text"
        />
        <button
          className="btn"
          onClick={() => {
            Register(name);
          }}
        >
          Register
        </button>
      </div>
    </div>
  );
}

export default App;
