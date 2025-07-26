import { createHashRouter } from "react-router";
import App from "./App";
import Chat from "./pages/Chat";
import UserRegister from "./pages/UserRegister";
import UserLogin from "./pages/UserLogin";

export const router = createHashRouter([
  {
    path: "/",
    Component: App,
    index: true,
  },
  {
    path: "/register",
    Component: UserRegister,
  },
  {
    path: "/login",
    Component: UserLogin,
  },
  {
    path: "/chat",
    Component: Chat,
  },
]);
