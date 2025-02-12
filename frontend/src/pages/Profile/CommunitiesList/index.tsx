import React from "react";
import arrasa from "assets/communities/arrasa.png";
import borracha from "assets/communities/borracha.png";
import campainha from "assets/communities/campainha.png";
import cinco from "assets/communities/cinco.png";
import duplas from "assets/communities/duplas.png";
import gina from "assets/communities/gina.png";
import odeio_acordar_cedo from "assets/communities/odeio_acordar_cedo.png";
import sorvete from "assets/communities/sorvete.png";
import tony from "assets/communities/tony.png";
import { palette } from "theme";

import { Box, Grid2 as Grid, Paper, Stack, Typography } from "@mui/material";

interface ICommunitiesListProps {
  userId: number;
}
export const CommunitiesList: React.FC<ICommunitiesListProps> = ({ userId }) => {
  const user = {
    friendsNumber: 132,
    communitiesNumber: 821,
  };

  const communities = [
    { id: 1, name: "Deus me disse: desce e arrasa", picture: arrasa, members: 450 },
    { id: 2, name: "Eu nunca terminei uma borracha", picture: borracha, members: 3214 },
    { id: 3, name: "Tocava campainha e corria", picture: campainha, members: 54352 },
    { id: 4, name: "Só mais 5 minutinhos", picture: cinco, members: 123 },
    { id: 5, name: "Não sei individualizar duplas", picture: duplas, members: 12 },
    { id: 6, name: "Tive medo da Gina dos palitos", picture: gina, members: 9 },
    { id: 7, name: "Eu odeio acordar cedo", picture: odeio_acordar_cedo, members: 839129 },
    { id: 8, name: "Queria sorvete, mas era feijão", picture: sorvete, members: 7 },
    { id: 9, name: "Eu colho flores e o Tony Ramos", picture: tony, members: 455 },
  ];

  return (
    <Paper>
      <Box sx={{ textAlign: "center", width: "100%", backgroundColor: palette.lightBlue }}>
        <Typography variant={"subtitle2"} fontWeight={"bold"}>
          {`comunidades (${user.communitiesNumber})`}
        </Typography>
      </Box>
      <Grid container spacing={1} sx={{ p: 1 }}>
        {communities.map((community) => (
          <Grid size={4} key={community.id}>
            <Stack alignItems={"center"} flexDirection={"column"} sx={{ textAlign: "center" }}>
              <img height={64} src={community.picture} alt="image not found" />
              <Typography variant={"body2"}>{`${community.name} (${community.members})`}</Typography>
            </Stack>
          </Grid>
        ))}
      </Grid>
    </Paper>
  );
};
