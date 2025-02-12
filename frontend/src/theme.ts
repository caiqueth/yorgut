import { ThemeOptions } from "@mui/material";

const lightBlue = "#eef4fb";
const mediumBlue = "#d5def0";
const darkBlue = "#cdd9ef";
const strongPink = "#d744b6";
const lightPink = "#f0cfe8";

export const palette = {
  lightBlue,
  mediumBlue,
  darkBlue,
  strongPink,
  lightPink,
};

export const theme: ThemeOptions = {
  palette: {
    primary: {
      main: strongPink,
    },
    secondary: {
      main: lightPink,
    },
    background: {
      paper: darkBlue,
      default: mediumBlue,
    },
  },
  typography: {
    fontFamily: "Verdana",
  },
  components: {
    MuiPaper: {
      styleOverrides: {
        root: {
          border: `1px solid ${darkBlue}`,
          boxShadow: "none",
          borderRadius: 0,
        },
      },
    },
  },
};
