import React, { useEffect, useState } from "react";
import { Link, useLocation } from "react-router-dom";

import { Box, Collapse, Typography, useTheme } from "@mui/material";
import { LinkBar } from "components/LinkBar";
import { LogoutLink } from "components/LogoutLink";
import { UserBar } from "components/UserBar";
import { UserEmail } from "components/UserEmail";
import { YorgutTitle } from "components/YorgutTitle";

// The least height of the bar on phones, so the links are easy to tap and the
// header's logo (Wrapper) lines up with them.
export const PHONE_BAR_HEIGHT = 28;

const APP_LINKS = [
  { to: "/", label: "Home" },
  { to: "/friends", label: "Amigos" },
  { to: "/scraps", label: "Mensagens" },
  { to: "/communities", label: "Comunidades" },
  { to: "/search", label: "Pesquisar" },
  { to: "/media", label: "Mídia" },
  { to: "/whats-new", label: "O que há de novo" },
];

// Shown on the phone bar next to "Mais" (more break the line on narrow
// phones), so the panel leaves them out.
const PHONE_BAR_LINKS = ["/search"];

const USER_LINKS = [
  { to: "/my-account", label: "Minha conta" },
  { to: "/help", label: "Ajuda" },
];

export const Navbar: React.FC = () => {
  const theme = useTheme();
  const { navbar, navbarEdge, navbarLogo } = theme.palette.yorgut;
  const { navbarPadding, navbarSeparators, logoInNavbar, navbarFullWidth, userBarInNavbar } = theme.yorgut;
  // background, not backgroundColor: the theme's fill may be a gradient. Inset
  // shadows, not borders, so the edges don't add to the height.
  // On phones, a theme with its logo outside the bar (Wrapper puts it beside
  // the bar there) leaves the bar unfilled: Wrapper fills the whole header,
  // logo included.
  const phoneFill = (value: string) => (logoInNavbar ? value : { xs: "none", sm: value });
  const fill = {
    background: phoneFill(navbar),
    boxShadow: phoneFill(`inset 0 1px ${navbarEdge}, inset 0 -1px ${navbarEdge}`),
  };

  // On phones the bar shows "Mais", a few links and "Sair"; "Mais" opens a
  // panel under it with the other links. Going to another page closes it.
  const [menuOpen, setMenuOpen] = useState(false);
  const { pathname } = useLocation();
  useEffect(() => setMenuOpen(false), [pathname]);

  const toggleMenu = (e: React.MouseEvent) => {
    e.preventDefault();
    setMenuOpen(open => !open);
  };

  const panelLink = { display: "block", py: 0.75 };

  return (
    <Box sx={{
      width: "100%",
      display: "flex",
      flexDirection: "column",
      // Full width: the fill goes on a layer as wide as the screen, centered
      // on the bar and behind its content (isolation keeps it above the page).
      ...(navbarFullWidth
        ? {
          position: "relative",
          isolation: "isolate",
          "&::before": {
            content: '""',
            position: "absolute",
            top: 0,
            bottom: 0,
            left: "50%",
            width: "100vw",
            transform: "translateX(-50%)",
            zIndex: -1,
            ...fill,
          },
        }
        : fill),
      color: "yorgut.navbarText",
      "& a, & a:visited": { color: "yorgut.navbarLink" },
      "& strong": { color: "yorgut.navbarEmail" },
    }}>
      <Box sx={{ display: "flex", alignItems: "stretch" }}>
        {/* Some themes put the logo here, in a box as tall as the bar; the
            others show it at the right of the header (Wrapper). */}
        {logoInNavbar && (
          <Box sx={{ display: "flex", alignItems: "center", flexShrink: 0, px: 1, background: navbarLogo }}>
            <YorgutTitle size={22} />
          </Box>
        )}
        <Box sx={{
          flex: 1,
          minWidth: 0,
          display: "flex",
          alignItems: "center",
          minHeight: { xs: PHONE_BAR_HEIGHT, sm: 0 },
          px: navbarPadding.x,
          py: navbarPadding.y,
        }}>
          {/* The app's links on the left, the user's on the right (when the
              theme puts them here); on narrow screens the user's go under. */}
          <Box sx={{ width: "100%", display: { xs: "none", sm: "flex" }, flexWrap: "wrap", alignItems: "center", columnGap: 2 }}>
            <Box sx={{ flex: "1 1 auto", minWidth: 0 }}>
              <LinkBar
                separators={navbarSeparators}
                items={APP_LINKS.map(({ to, label }) => <Link to={to}>{label}</Link>)}
              />
            </Box>
            {userBarInNavbar && (
              <Box sx={{ flex: "0 1 auto", minWidth: 0 }}>
                <UserBar separators={navbarSeparators} />
              </Box>
            )}
          </Box>
          <Typography
            variant={"body2"}
            component={"div"}
            sx={{ width: "100%", display: { xs: "flex", sm: "none" }, justifyContent: "space-between", alignItems: "center", columnGap: 2 }}
          >
            <Box sx={{ minWidth: 0 }}>
              <LinkBar
                separators={navbarSeparators}
                items={[
                  <a href="#" role="button" aria-expanded={menuOpen} onClick={toggleMenu}>
                    Mais {menuOpen ? "▴" : "▾"}
                  </a>,
                  ...APP_LINKS
                    .filter(({ to }) => PHONE_BAR_LINKS.includes(to))
                    .map(({ to, label }) => <Link to={to}>{label}</Link>),
                ]}
              />
            </Box>
            <LogoutLink />
          </Typography>
        </Box>
      </Box>
      <Collapse in={menuOpen} sx={{ display: { sm: "none" } }}>
        <Typography variant={"body2"} component={"nav"} sx={{ px: Math.max(navbarPadding.x, 1), pb: 1 }}>
          {APP_LINKS.filter(({ to }) => !PHONE_BAR_LINKS.includes(to)).map(({ to, label }) => (
            <Link key={to} to={to}><Box component={"span"} sx={panelLink}>{label}</Box></Link>
          ))}
          <Box sx={{ my: 0.5, borderTop: "1px solid", opacity: 0.3 }} />
          <Box sx={panelLink}><UserEmail /></Box>
          {USER_LINKS.map(({ to, label }) => (
            <Link key={to} to={to}><Box component={"span"} sx={panelLink}>{label}</Box></Link>
          ))}
        </Typography>
      </Collapse>
    </Box>
  );
};
