import React from "react";
import { Wrapper } from "components/Wrapper";
import { Profile } from "pages/Profile";

import { CssBaseline } from "@mui/material";

import "./App.css";

const App: React.FC = () => {
  return (
    <CssBaseline>
      <Wrapper>
        <Profile />
      </Wrapper>
    </CssBaseline>
  );
};

export default App;
