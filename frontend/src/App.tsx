import React, { useMemo, useState } from "react";
import { Wrapper } from "components/Wrapper";
import { useNavigate, useRoutes } from "react-router-dom";
import { Routes } from "routes";
import { ApiError } from "api/http";

import { Alert, Snackbar } from "@mui/material";
import {
  MutationCache,
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";

import { AppThemeProvider } from "components/AppThemeProvider";

import "App.css";
import { GlobalContextStateProvider } from "contexts/GlobalStateContext";
import { AuthContextProvider } from "contexts/AuthContext";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

const App: React.FC = () => {
  const navigate = useNavigate();
  const routes = useRoutes(Routes);

  const queryCache = useMemo(
    () =>
      new QueryCache({
        onError: (error) => {
          if (!(error instanceof ApiError)) return;
          if (error.status === 401) navigate("/login");
          else if (error.status === 404) navigate("/not-found");
        },
      }),
    [],
  );

  // The message of the last failed action (sending a scrap, joining a
  // community...), shown at the bottom of the screen. Forms that show their
  // own errors mark their mutations with meta.showsErrorInline.
  const [actionError, setActionError] = useState<string | null>(null);
  const mutationCache = useMemo(
    () =>
      new MutationCache({
        onError: (error, _variables, _context, mutation) => {
          if (error instanceof ApiError && error.status === 401) {
            navigate("/login");
            return;
          }
          if (mutation.meta?.showsErrorInline) return;
          setActionError(error.message);
        },
      }),
    [],
  );

  const queryClient = useMemo(
    () =>
      new QueryClient({
        queryCache,
        mutationCache,
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
      <AppThemeProvider>
        <AuthContextProvider>
          <GlobalContextStateProvider>
            <Wrapper>
              {routes}
            </Wrapper>
            <Snackbar
              open={!!actionError}
              autoHideDuration={6000}
              onClose={() => setActionError(null)}
              anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
            >
              <Alert severity={"error"} variant={"filled"} onClose={() => setActionError(null)}>
                {actionError}
              </Alert>
            </Snackbar>
          </GlobalContextStateProvider>
        </AuthContextProvider>
      </AppThemeProvider>
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
};

export default App;
