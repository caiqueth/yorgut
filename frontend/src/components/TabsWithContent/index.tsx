import { Box, Stack } from "@mui/material";
import { Button } from "components/Button";

interface IProps {
  tabs: {
    title: string,
    onClick: () => void
  }[];
  extraHeader: React.ReactNode;
  children?: React.ReactNode;
}
// Tabs over the content, with extraHeader (e.g. a "+ criar" button) at their
// right. On phones extraHeader moves above the tabs, and the tabs keep their
// text on one line, scrolling sideways if they still don't fit.
export const TabsWithContent: React.FC<IProps> = ({ tabs, extraHeader, children }) => {
  return (
    <Box sx={{ width: "100%" }}>
      <Stack direction={{ xs: "column-reverse", sm: "row" }} alignItems={{ xs: "stretch", sm: "end" }} gap={{ xs: 1, sm: 0 }}>
        <Stack direction={"row"} spacing={"1px"} alignItems={"end"} sx={{ overflowX: "auto", minWidth: 0 }}>
          {tabs.map(({ title, onClick }, idx) => (
            <Button sx={{
              flexShrink: 0,
              whiteSpace: "nowrap",
              px: { xs: 1, sm: 2 },
              borderRadius: "4px 4px 0 0",
              boxShadow: "none",
              borderTop: 2,
              borderLeft: 2,
              borderRight: 2,
              borderColor: "yorgut.border",
              // background, so it covers Button's (which may be a gradient).
              background: (theme) => theme.palette.yorgut.tab,
              color: "yorgut.tabText",
            }}
              key={idx}
              onClick={onClick}>
              {title}
            </Button>
          ))}
        </Stack>
        <Box flexGrow={1} />
        <Box sx={{ alignSelf: { xs: "end", sm: "auto" } }}>{extraHeader}</Box>
      </Stack>
      {children}
    </Box>
  )
}
