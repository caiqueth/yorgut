import React, { useMemo } from "react";
import { Wrapper } from "components/templates/Wrapper";
import { useNavigate, useRoutes } from "react-router-dom";
import { Routes } from "routes";

import { createTheme, CssBaseline, ThemeProvider } from "@mui/material";
import {
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";

import { theme } from "./theme";

import "./App.css";

const App: React.FC = () => {
  const navigate = useNavigate();
  const routes = useRoutes(Routes);

  const queryCache = useMemo(
    () =>
      new QueryCache({
        onError: (error) => {
          if (error.message.toLowerCase().startsWith("unauthorized")) {
            navigate("/login");
            return;
          }
          if (error.message.toLowerCase().startsWith("not allowed")) {
            navigate("/not-allowed");
            return;
          }
          if (error.message.toLowerCase().startsWith("not found")) {
            navigate("/not-found");
            return;
          }
          // navigate("/bad-server");
          // return;
        },
      }),
    [],
  );

  const queryClient = useMemo(
    () =>
      new QueryClient({
        queryCache,
        defaultOptions: {
          queries: {
            retry: false,
          },
        },
      }),
    [],
  );

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={createTheme(theme)}>
        <CssBaseline />
        <Wrapper>{routes}</Wrapper>
      </ThemeProvider>
    </QueryClientProvider>
  );
};

export default App;
