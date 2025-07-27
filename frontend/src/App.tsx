import { useEffect } from "react";
import { GetProfile } from "../wailsjs/go/user/UserService";
import type { TProfileSchema, TResponseSchema } from "@/models";
import { useNavigate } from "react-router";
import MainLayout from "@/components/MainLayout";
import { Loader } from "@/components/Loader";

function App() {
  const navigate = useNavigate();

  useEffect(() => {
    GetProfile().then((res: TResponseSchema<TProfileSchema>) => {
      if (res.code === 200) {
        if (res.data.id) {
          navigate("/login");
          return;
        }
      }

      navigate("/register");
    });
  }, []);

  return (
    <MainLayout>
      <Loader />
    </MainLayout>
  );
}

export default App;
