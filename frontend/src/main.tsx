import { createRoot } from "react-dom/client";
import "./style.css";
import { RouterProvider } from "react-router";
import { router } from "./routes";

const container = document.getElementById("root");

const root = createRoot(container!);

root.render(<RouterProvider router={router} />);
