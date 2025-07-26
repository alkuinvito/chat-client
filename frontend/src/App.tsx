import MainLayout from "./components/MainLayout";
import { useState, useEffect } from "react";
import { EventsOn } from "../wailsjs/runtime/runtime";
import { Register } from "../wailsjs/go/auth/AuthService";
import { useNavigate } from "react-router";
import { Input } from "./components/ui/input";
import { Button } from "./components/ui/button";
import { Label } from "./components/ui/label";

function App() {
  const [name, setName] = useState("");

  const navigate = useNavigate();

  useEffect(() => {
    EventsOn("auth:authorized", (path) => {
      navigate(path);
    });
  }, [navigate]);

  return (
    <MainLayout className="flex flex-col justify-center gap-16">
      <div>
        <h1 className="text-4xl font-bold">P2P Chat Client</h1>
      </div>
      <div className="grid gap-4 w-full max-w-[320px] mx-auto">
        <Label htmlFor="username">Username</Label>
        <Input
          id="username"
          name="username"
          onChange={(e) => {
            setName(e.target.value);
          }}
        />
        <Button
          variant="secondary"
          onClick={() => {
            Register(name);
          }}
        >
          Register
        </Button>
      </div>
    </MainLayout>
  );
}

export default App;
