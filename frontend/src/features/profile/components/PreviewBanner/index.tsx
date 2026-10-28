import React from "react";
import { Alert, Box } from "@mui/material";
import { Link } from "react-router-dom";
import { Button } from "components/Button";
import { useProfilePreview } from "contexts/ProfilePreviewContext";

// Over the profile while "ver como visitante" is on.
export const PreviewBanner: React.FC = () => {
  const { previewing, setPreviewing } = useProfilePreview();
  if (!previewing) return null;

  return (
    <Box sx={{ mb: 2 }}>
      <Alert
        severity={"info"}
        action={<Button size={"small"} onClick={() => setPreviewing(false)}>sair da prévia</Button>}
      >
        Você está vendo seu perfil como quem não é seu amigo vê.{" "}
        Mude o que aparece em <Link to="/my-account">minha conta</Link>.
      </Alert>
    </Box>
  );
};
