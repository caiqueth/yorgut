import { Box, Stack, Typography } from "@mui/material";
import { useLoggedUser } from "api/user/queries";
import { Link, useNavigate } from "react-router-dom";
import { googleLogout } from "@react-oauth/google";

export const UserBar: React.FC = () => {
  const { email } = useLoggedUser();

  const navigate = useNavigate();

  const handleLogout = () => {
    googleLogout();
    navigate("/login");
  };

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
      <Typography variant={"body2"}>
        <Link to="#" onClick={() => handleLogout()}>
          Sair
        </Link>
      </Typography>
    </Stack>
  );
};
