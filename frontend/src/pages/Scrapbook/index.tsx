import {
  Box,
  Grid2 as Grid,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useLoggedUser } from "api/user/queries";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { UserProfilePicture } from "components/organisms/UserProfilePicture";
import { Menu } from "pages/Home/components/Menu";
import React from "react";

interface IProps {}
export const Scrapbook: React.FC<IProps> = () => {
  const { userId } = useLoggedUser();

  const items: React.ReactNode[] = [
    <Typography>Teste de recado 1</Typography>,
    <Typography>Teste de recado 2</Typography>,
  ];

  return (
    <Box sx={{ width: "100%" }}>
      <Grid container spacing={2}>
        <Grid size={2}>
          <Stack direction={"column"} gap={2}>
            <UserProfilePicture userId={userId} />
            <Menu userId={userId} />
          </Stack>
        </Grid>
        <Grid size={10}>
          <Stack direction={"column"} gap={1}>
            <Box sx={{ width: "100%" }}>
              <Typography>Deixar um recado</Typography>
              <TextField multiline rows={6} fullWidth></TextField>
            </Box>
            <BoxWithTitle title="meus recados">
              <StripedGrid items={items} elementsPerRow={1} />
            </BoxWithTitle>
          </Stack>
        </Grid>
      </Grid>
    </Box>
  );
};
