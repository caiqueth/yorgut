import { createTheme, Theme, ThemeOptions } from "@mui/material";
import { FindMeResponseTheme } from "api/generated/model";

// The app's own colors, by what they paint, so a theme can change one place
// without the others (the navbar but not the lists, say). Components read
// them as theme.palette.yorgut.<role>, never as hex values.
export interface YorgutColors {
  // Box borders (every Paper), tab outlines and the lines around the header
  // and footer.
  border: string;
  // The navbar's fill: a color or a CSS gradient. Painted as a background, so
  // a gradient works.
  navbar: string;
  // A 1px line at the navbar's top and bottom edges, drawn inside it
  // (transparent for none).
  navbarEdge: string;
  // The navbar's separators and its links.
  navbarText: string;
  navbarLink: string;
  // The user's e-mail, when the navbar holds the user's links
  // (yorgut.userBarInNavbar).
  navbarEmail: string;
  // Behind the logo, when the navbar holds it (yorgut.logoInNavbar): a color
  // or a CSS gradient.
  navbarLogo: string;
  // Every link outside the navbar, visited or not.
  link: string;
  // The strip behind a BoxWithTitle's title.
  boxTitle: string;
  // The light fill of side-menu items (LateralButtonList) and of the
  // formatting preview.
  highlight: string;
  // The text of side-menu items, and the lines between and around them
  // (transparent to let the box show through).
  menuText: string;
  menuDivider: string;
  // Button's fill: a color or a CSS gradient. Its glow inside the edges and
  // its 1px outline (transparent for none).
  button: string;
  buttonShadow: string;
  buttonBorder: string;
  buttonText: string;
  // The tabs over a TabsWithContent.
  tab: string;
  tabText: string;
  // Alternating rows of lists (StripedGrid, pending invites): striped is the
  // first row and every other one after it (helpers/isStripedRow).
  rows: { striped: string; plain: string };
  // Alternating rows of ProfileTable, first row striped, and the color of
  // its labels ("relacionamento:").
  tableRows: { striped: string; plain: string };
  tableLabel: string;
}

// What a theme changes besides colors, read as theme.yorgut.<setting>.
export interface YorgutLayout {
  // Where a BoxWithTitle's title sits.
  boxTitleAlign: "center" | "left";
  // The corners of a page's main box (the home's welcome, the profile's
  // header), which may stand out from the other boxes' (a CSS border-radius).
  featuredBoxRadius: string;
  // Space around the navbar's links, in theme spacing units: x on the sides,
  // y above and below.
  navbarPadding: { x: number; y: number };
  // Whether the navbar's links have "|" between them.
  navbarSeparators: boolean;
  // Whether the logo sits at the left of the navbar, in its own box, instead
  // of at the right of the header past a diagonal line.
  logoInNavbar: boolean;
  // Whether the navbar's fill runs to both edges of the screen; its links
  // stay in the page's column either way.
  navbarFullWidth: boolean;
  // Whether the header touches the top of the screen, with no space above.
  headerAtTop: boolean;
  // Whether the user's links (e-mail, Minha conta, Ajuda, Sair) sit at the
  // right of the navbar instead of on their own line above it.
  userBarInNavbar: boolean;
  // When links are underlined: always, or only under the pointer.
  linkUnderline: "always" | "hover";
  // Space between side-menu items (LateralButtonList), in pixels.
  menuGap: number;
  // White space between ProfileTable's rows and around them, in pixels.
  tableRowSpacing: number;
  // Whether the profile's header and its fields make one box, the header
  // running into the table, instead of two boxes apart.
  joinedProfileHeader: boolean;
  // Button's size and text, instead of MUI's (for the button's size) with
  // 16px on the sides.
  buttonShape?: { padding: string; fontSize: string; fontWeight: number };
}

declare module "@mui/material/styles" {
  interface Palette {
    yorgut: YorgutColors;
  }
  interface PaletteOptions {
    yorgut?: Partial<YorgutColors>;
  }
  interface Theme {
    yorgut: YorgutLayout;
  }
  interface ThemeOptions {
    yorgut?: Partial<YorgutLayout>;
  }
}

const lightBlue = "#eef4fb";
const mediumBlue = "#d5def0";
const darkBlue = "#cdd9ef";
const darkerBlue = "#bac6dd";
const strongPink = "#d744b6";
const lightPink = "#f0cfe8";

// The 2005 look, the app's first and default one. Other themes are changes on
// top of it.
const theme2005: ThemeOptions = {
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
    yorgut: {
      border: darkerBlue,
      navbar: darkBlue,
      navbarEdge: "transparent",
      navbarText: "rgba(0, 0, 0, 0.87)",
      navbarLink: "LinkText",
      navbarEmail: "inherit",
      navbarLogo: "transparent",
      // The browser's own link blue.
      link: "LinkText",
      boxTitle: lightBlue,
      highlight: lightBlue,
      menuText: "rgba(0, 0, 0, 0.87)",
      menuDivider: "transparent",
      button: mediumBlue,
      buttonShadow: "rgb(146, 163, 195)",
      buttonBorder: "transparent",
      buttonText: "black",
      tab: mediumBlue,
      tabText: "black",
      rows: { striped: darkBlue, plain: mediumBlue },
      tableRows: { striped: lightBlue, plain: "white" },
      tableLabel: "rgba(0, 0, 0, 0.87)",
    },
  },
  yorgut: {
    boxTitleAlign: "center",
    featuredBoxRadius: "0",
    navbarPadding: { x: 0, y: 0 },
    navbarSeparators: true,
    logoInNavbar: false,
    navbarFullWidth: false,
    headerAtTop: false,
    userBarInNavbar: false,
    linkUnderline: "always",
    menuGap: 2,
    tableRowSpacing: 0,
    joinedProfileHeader: false,
  },
  typography: {
    // Verdana is installed on Windows and macOS (the app doesn't ship it: its
    // license forbids redistributing the file); the rest are look-alikes.
    // "Color Emoji" is defined in index.css.
    fontFamily: '"Color Emoji", Verdana, Tahoma, "DejaVu Sans", Geneva, Arial, sans-serif',
    // Verdana runs large, so every variant is scaled down from MUI's default 14
    // (body2 becomes 12px, body1 about 14px).
    fontSize: 12,
  },
  components: {
    // Links keep their color after being visited, instead of turning purple.
    MuiCssBaseline: {
      styleOverrides: (theme) => ({
        "a, a:visited": {
          color: theme.palette.yorgut.link,
          textDecoration: theme.yorgut.linkUnderline === "always" ? "underline" : "none",
        },
        "a:hover": { textDecoration: "underline" },
        // Full-width pieces inside the page's column (the navbar's fill)
        // are 100vw wide, which counts the scrollbar: clip, so the page
        // doesn't scroll sideways by its width.
        body: { overflowX: "clip" },
      }),
    },
    MuiPaper: {
      styleOverrides: {
        // A function, so a theme that changes only the border color still
        // recolors every box.
        root: ({ theme }) => ({
          border: `2px solid ${theme.palette.yorgut.border}`,
          boxShadow: "none",
          borderRadius: 0,
        }),
      },
    },
    // The light fill goes on the input itself, not on the TextField: the
    // TextField's box also holds the floating label and the helper text, so a
    // fill there has square corners past the rounded border, cuts the label
    // in two and paints behind the helper text.
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          backgroundColor: "#fafafa",
        },
      },
    },
    // MUI sizes fields for 16px text (body1, 56px tall), but the app's text is
    // mostly body2 (14px), so default fields looked oversized next to it.
    MuiTextField: {
      defaultProps: {
        size: "small",
        // Every field's look; a field that needs another passes variant itself.

      },
    },
    MuiInputBase: {
      styleOverrides: {
        root: ({ theme }) => ({ fontSize: theme.typography.body2.fontSize }),
      },
    },
    MuiInputLabel: {
      styleOverrides: {
        root: ({ theme }) => ({ fontSize: theme.typography.body2.fontSize }),
      },
    },
    // Rectangular by default (pictures, boxes). A skeleton standing for text
    // must pass variant="text": only that one takes the height of its line,
    // a rectangular one without a height is 0px tall.
    MuiSkeleton: {
      defaultProps: {
        variant: "rectangular",
        animation: "wave",
      }
    }
  },
};

// The 2008 redesign: white boxes with rounded, uneven corners on a light blue
// page, a thicker gradient bar for the links, grayer text, and titles on the
// left. Colors sampled from screenshots of the time.
const theme2008: ThemeOptions = {
  palette: {
    background: {
      paper: "#ffffff",
      default: "#dae5f7",
    },
    text: {
      primary: "#222222",
      secondary: "#808488",
    },
    yorgut: {
      border: "#c9d7ea",
      // Lightest a little under the middle, darker at both edges.
      navbar: "linear-gradient(to bottom, #5b88c1 0%, #77a0ce 65%, #5b8ac2 100%)",
      navbarEdge: "#4c6f9d",
      navbarText: "#ffffff",
      navbarLink: "#ffffff",
      navbarEmail: "#234475",
      navbarLogo: "linear-gradient(to bottom, #fbfbfb 0%, #f4f5fa 50%, #d4d5da 80%, #e2e3e7 100%)",
      link: "#2b6176",
      // Side-menu items read as links, as they are.
      menuText: "#2b6176",
      menuDivider: "#dfe9f3",
      boxTitle: "#ffffff",
      highlight: "#eff8ff",
      // Small white pills, a little grayer at the bottom.
      button: "linear-gradient(to bottom, #ffffff 55%, #eef0f3 100%)",
      buttonShadow: "transparent",
      buttonBorder: "#c8d1dd",
      buttonText: "#55779c",
      tab: "#7d98b5",
      tabText: "#ffffff",
      rows: { striped: "#f1f1f1", plain: "#f8f8f8" },
      tableRows: { striped: "#dfeeff", plain: "#eef9ff" },
      tableLabel: "#66717b",
    },
  },
  yorgut: {
    boxTitleAlign: "left",
    // Much rounder on the top right, the redesign's uneven corner.
    featuredBoxRadius: "8px 40px 8px 8px",
    navbarPadding: { x: 1.5, y: 1 },
    navbarSeparators: false,
    logoInNavbar: true,
    navbarFullWidth: true,
    headerAtTop: true,
    userBarInNavbar: true,
    linkUnderline: "hover",
    menuGap: 1,
    tableRowSpacing: 2,
    joinedProfileHeader: true,
    buttonShape: { padding: "4px 10px", fontSize: "11px", fontWeight: 700 },
  },
  components: {
    MuiPaper: {
      styleOverrides: {
        // Rounded and lifted a little off the page. A page's main box is
        // rounder on one corner (yorgut.featuredBoxRadius).
        root: ({ theme }) => ({
          border: `1px solid ${theme.palette.yorgut.border}`,
          boxShadow: "0 1px 3px rgba(60, 90, 140, 0.18)",
          borderRadius: 8,
        }),
      },
    },
  },
};

// The looks a user may pick, by the key the API stores (FindMeResponse.theme).
export type ThemeKey = FindMeResponseTheme;

export const defaultTheme: ThemeKey = "2005";

// createTheme deep-merges its arguments, so a theme only lists what it changes.
export const themes: Record<ThemeKey, Theme> = {
  "2005": createTheme(theme2005),
  "2008": createTheme(theme2005, theme2008),
};

// The picker's options, in the order it shows them.
export const themeOptions: { key: ThemeKey; label: string; description: string }[] = [
  { key: "2005", label: "clássico (2005)", description: "O visual do começo: cantos retos e azul-claro." },
  { key: "2008", label: "novo (2008)", description: "O visual de 2008: cantos arredondados e cores novas." },
];

export const isThemeKey = (value: unknown): value is ThemeKey =>
  typeof value === "string" && value in themes;
