import { CssBaseline, ThemeProvider } from "@mui/material";
import { useMe } from "api/me/queries";
import React, { useEffect } from "react";
import { defaultTheme, isThemeKey, themes } from "theme";

// The last theme of this browser, so a reload starts on it instead of
// flashing the default one until /me answers. The saved one is the user's
// (FindMeResponse.theme); this is only a copy.
const storageKey = "theme";

const readCachedTheme = () => {
  try {
    const cached = localStorage.getItem(storageKey);
    return isThemeKey(cached) ? cached : defaultTheme;
  } catch {
    return defaultTheme;
  }
};

// Applies the logged user's theme to the whole app. Logged out, the last
// theme used in this browser.
export const AppThemeProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
  const { data: me } = useMe();
  // A theme this build doesn't know (renamed or removed) falls back to the
  // default instead of leaving the app without one.
  const key = me ? (isThemeKey(me.theme) ? me.theme : defaultTheme) : readCachedTheme();

  useEffect(() => {
    if (!me || !isThemeKey(me.theme)) return;
    try {
      localStorage.setItem(storageKey, me.theme);
    } catch {
      // Not cached: the next reload starts on the default theme.
    }
  }, [me?.theme]);

  return (
    <ThemeProvider theme={themes[key]}>
      <CssBaseline />
      {children}
    </ThemeProvider>
  );
};
