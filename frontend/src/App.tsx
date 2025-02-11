import React from "react";
import { Wrapper } from "components/Wrapper";
import { Profile } from "pages/Profile";

import { createTheme, CssBaseline, ThemeProvider } from "@mui/material";

import { theme } from "./theme";

import "./App.css";

const App: React.FC = () => {
  return (
    <ThemeProvider theme={createTheme(theme)}>
      <CssBaseline />
      <Wrapper>
        <Profile />
      </Wrapper>
    </ThemeProvider>
  );
};

export default App;
