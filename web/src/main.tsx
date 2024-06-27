import * as React from "react";
import * as ReactDOM from "react-dom/client";
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import "./index.css";
import { Container } from "react-dom/client";
import Root from "./routes/root.tsx";
import Login from "./routes/login.tsx";
import { GoogleOAuthProvider } from "@react-oauth/google";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
  },
  {
    path: "/login",
    element: <Login />,
  },
]);

ReactDOM.createRoot(document.getElementById("root") as Container).render(
  <React.StrictMode>
    <GoogleOAuthProvider clientId="962765624064-p4qkev7rg3jo2tt0va47j8t24mmhqd7u.apps.googleusercontent.com">
      <RouterProvider router={router} />
    </GoogleOAuthProvider>
  </React.StrictMode>
);
