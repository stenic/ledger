import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App";
import { BrowserRouter } from "react-router-dom";
// import reportWebVitals from "./reportWebVitals";
import { theme } from "./theme";
import { CssBaseline, ThemeProvider } from "@mui/material";
import { AuthProvider } from "react-oidc-context";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement
);

fetch("/auth/config")
  .then((e) => e.json())
  .then((d) => {
    const oidcConfig = {
      ...d,
      redirect_uri: window.location.href,
      onSigninCallback: () => {
        window.history.replaceState(
          {},
          document.title,
          window.location.pathname
        );
      },
    };

    const queryClient = new QueryClient();

    root.render(
      <React.StrictMode>
        <QueryClientProvider client={queryClient}>
          <BrowserRouter>
            <ThemeProvider theme={theme}>
              <CssBaseline />
              <AuthProvider {...oidcConfig}>
                <App />
              </AuthProvider>
            </ThemeProvider>
          </BrowserRouter>
          <ReactQueryDevtools initialIsOpen={false} />
        </QueryClientProvider>
      </React.StrictMode>
    );
  });

// reportWebVitals();
