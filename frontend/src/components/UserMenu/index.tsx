import { Stack, Typography } from "@mui/material";

export const UserMenu: React.FC = () => {
  return (
    <Stack direction={"row"} gap={1} sx={{ width: "100%", justifyContent: "end" }}>
      <Typography variant={"body2"}>
        <strong>email@provedor.com</strong>
      </Typography>
      <Typography variant={"body2"}>|</Typography>
      <Typography variant={"body2"}>Minha conta</Typography>
      <Typography variant={"body2"}>|</Typography>
      <Typography variant={"body2"}>Ajuda</Typography>
      <Typography variant={"body2"}>|</Typography>
      <Typography variant={"body2"}>Sair</Typography>
    </Stack>
  );
};
