import React from "react";
import { palette } from "theme";

import { Divider, List, ListItem, ListItemAvatar, ListItemButton, ListItemText, Paper } from "@mui/material";

export const ProfileMenu: React.FC = () => {
  return (
    <Paper>
      <List sx={{ backgroundColor: palette.lightBlue }}>
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>📝</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"recados"} />
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
            <ListItemAvatar>📹</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"vídeos"} />
          </ListItemButton>
        </ListItem>
        <Divider />
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>📧</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"enviar mensagem"} />
          </ListItemButton>
        </ListItem>
        <Divider />
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>✒</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"escrever depoimento"} />
          </ListItemButton>
        </ListItem>
        <Divider />
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>😋</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"provocar"} />
          </ListItemButton>
        </ListItem>
        <Divider />
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>🤐</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"ignorar"} />
          </ListItemButton>
        </ListItem>
        <Divider />
        <ListItem disablePadding>
          <ListItemButton>
            <ListItemAvatar>🚫</ListItemAvatar>
            <ListItemText slotProps={{ primary: { fontSize: "10px" } }} primary={"reportar"} />
          </ListItemButton>
        </ListItem>
      </List>
    </Paper>
  );
};
