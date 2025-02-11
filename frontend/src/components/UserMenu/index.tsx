import { Stack, Typography } from "@mui/material";

export const UserMenu: React.FC = () => {
  return (
    <Stack>
      <Typography>email@provedor.com</Typography>
      <Typography>Minha conta</Typography>
      <Typography>Ajuda</Typography>
      <Typography>Sair</Typography>
    </Stack>
  );
};
