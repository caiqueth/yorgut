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
    <GoogleOAuthProvider clientId={__GOOGLE_CLIENT_ID__}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </GoogleOAuthProvider>
  </React.StrictMode>,
);
