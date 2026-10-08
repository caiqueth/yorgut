import React from "react";

import { Box, Paper, Skeleton, Stack, Typography } from "@mui/material";
import { dayjsFromDate } from "helpers/formatting/date";
import { useMeWelcome } from "api/me/queries";
import { Link } from "react-router-dom";

interface IProps {

}
export const Welcome: React.FC<IProps> = () => {
  const { data } = useMeWelcome();

  return (
    <Box>
      <Paper sx={{ p: 1 }}>
        {data
          ? (
            <Stack direction={"column"} gap={1}>
              <Typography
                variant={"subtitle1"}
                sx={{ "&::first-letter": { fontSize: 36, lineHeight: 1 } }}
              >
                Bem-vindo(a), {data.userFirstName}
              </Typography>
              <Box>
                <Typography variant={"body1"}>
                  Você está conectado(a) a {data.indirectConnections}{" "}
                  pessoas através de {data.directConnections} amigos
                </Typography>
                <Typography variant={"body1"}>
                  Aumente sua rede agora mesmo!
                </Typography>
              </Box>
              <Typography variant={"body1"}>
                ⭐ Você tem {data.numberOfFans} fãs.
              </Typography>
              <Typography variant={"body1"}>
                📝 {data.numberOfScraps} recados.
              </Typography>
              <Typography variant={"body1"}>
                <strong>Visualizações do seu perfil:</strong>
              </Typography>
              <Typography variant={"body1"}>
                Desde {dayjsFromDate(data.userCreatedAt, "D [de] MMMM [de] YYYY")}:{" "}
                {data.profileViews.allTime}. Semana passada:{" "}
                {data.profileViews.lastWeek}.
              </Typography>
              <Box>
                <Typography variant={"body1"}>
                  <strong>Visitantes recentes:</strong>
                </Typography>
                <Typography variant={"body1"}>
                  {data.recentVisitors.length > 0 ? data.recentVisitors
                    .map(({ userId, firstName }) => (
                      <Link to={`/profile/${userId}`}>
                        <Typography>{firstName}</Typography>
                      </Link>
                    )) : "Por enquanto, ninguém."}
                </Typography>
              </Box>
              <Box>
                <Typography variant={"body1"}>
                  <strong>Próximas etapas:</strong>
                </Typography>
                <Typography variant={"body1"}>
                  <Link to="/friends">Classifique seus amigos!</Link>
                </Typography>
              </Box>
              <Box>
                <Typography variant={"body1"}>
                  <strong>Sorte de hoje</strong>
                </Typography>
                <Typography variant={"body1"}>{data.todaysLuck}</Typography>
              </Box>
            </Stack>
          )
          : <Skeleton height={392} />
        }
      </Paper>
    </Box>
  );
};
