import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";

import App from "./App";

import "./index.css";
import dayjs from "dayjs";
import { GoogleOAuthProvider } from "@react-oauth/google";

require("dayjs/locale/pt-br");

dayjs.locale("pt-br");

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement,
);
root.render(
  <React.StrictMode>
    <GoogleOAuthProvider
      clientId={
        "461443292884-v7mal5ohpg6srh7ubjhpphq2ol0vl7j6.apps.googleusercontent.com"
      }
    >
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </GoogleOAuthProvider>
  </React.StrictMode>,
);
