import { createHashRouter } from "react-router";
import App from "./App";
import Chat from "./pages/Chat";

export const router = createHashRouter([
  {
    path: "/",
    Component: App,
    index: true,
  },
  {
    path: "/chat",
    Component: Chat,
  },
]);
