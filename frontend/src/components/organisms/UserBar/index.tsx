import { Box, Stack, Typography } from "@mui/material";
import { useLoggedUser } from "api/user/queries";
import { Link } from "react-router-dom";

export const UserBar: React.FC = () => {
  const { email } = useLoggedUser();

  const handleLogout = () => {};

  return (
    <Stack
      direction={"row"}
      gap={1}
      sx={{ width: "100%", justifyContent: "end" }}
    >
      <Typography variant={"body2"}>
        <strong>{email}</strong>
      </Typography>
      <Typography variant={"body2"}>|</Typography>
      <Typography variant={"body2"}>
        <Link to="/my-account">Minha conta</Link>
      </Typography>
      <Typography variant={"body2"}>|</Typography>
      <Typography variant={"body2"}>
        <Link to="/help">Ajuda</Link>
      </Typography>
      <Typography variant={"body2"}>|</Typography>
      <Typography variant={"body2"}>Sair</Typography>
    </Stack>
  );
};
