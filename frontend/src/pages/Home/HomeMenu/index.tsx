import React from "react";
import { palette } from "theme";

import { Box, Divider, List, ListItem, ListItemAvatar, ListItemButton, ListItemText } from "@mui/material";

interface IHomeMenuProps {
  userId: number;
}
export const HomeMenu: React.FC<IHomeMenuProps> = ({ userId }) => {
  return (
    <Box>
      <List sx={{ backgroundColor: palette.lightBlue }}>
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>🙂</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"perfil"} />
          </ListItemButton>
        </ListItem>
        <Divider />
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>📷</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"álbum"} />
          </ListItemButton>
        </ListItem>
        <Divider />
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>📃</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"listas"} />
          </ListItemButton>
        </ListItem>
        <Divider />
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>📝</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"recados"} />
          </ListItemButton>
        </ListItem>
        <Divider />
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>✒</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"depoimentos"} />
          </ListItemButton>
        </ListItem>
        <Divider />
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>🔨</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"configurações"} />
          </ListItemButton>
        </ListItem>
      </List>
    </Box>
  );
};
