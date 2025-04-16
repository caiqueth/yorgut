import React, { useMemo } from "react";
import { Wrapper } from "components/Wrapper";
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
        onError: (error: Error) => {
          if (error.message.startsWith("Unauthorized")) {
            navigate("/disconnected");
          }
          if (error.message.startsWith("Not allowed")) {
            navigate("/not-allowed");
          }
        },
      }),
    [],
  );

  const queryClient = useMemo(
    () =>
      new QueryClient({
        queryCache,
        defaultOptions: { queries: { retry: false } },
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
