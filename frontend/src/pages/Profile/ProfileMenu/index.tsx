import React from "react";

import { List, ListItem } from "@mui/material";

export const ProfileMenu: React.FC = () => {
  return (
    <List>
      <ListItem>recados</ListItem>
      <ListItem>álbum</ListItem>
      <ListItem>vídeos</ListItem>
      <ListItem>enviar mensagem</ListItem>
      <ListItem>escrever depoimento</ListItem>
      <ListItem>provocar</ListItem>

      <ListItem>ignorar</ListItem>
      <ListItem>reportar</ListItem>
    </List>
  );
};
