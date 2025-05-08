import React from "react";
import { useTodaysLuck, useUserStatistics } from "api/user/queries";

import { Box, Paper, Stack, Typography } from "@mui/material";
import { dayjsFromDate } from "helpers/formatting/date";
import { useMe } from "api/me/queries";

interface IProps {

}
export const Welcome: React.FC<IProps> = () => {
  const { data: user } = useMe();

  if (!user) return null;

  const userStatistics = useUserStatistics();
  const todaysLuck = useTodaysLuck();

  return (
    <Box>
      <Paper sx={{ p: 1 }}>
        <Stack direction={"column"} gap={1}>
          <Typography
            variant={"subtitle1"}
            sx={{ "&::first-letter": { fontSize: "200%" } }}
          >
            Bem-vindo(a), {user.firstName}
          </Typography>
          <Box>
            <Typography variant={"body1"}>
              Você está conectado(a) a {userStatistics.undirectConnections}{" "}
              pessoas através de {userStatistics.totalFriends} amigos
            </Typography>
            <Typography variant={"body1"}>
              Aumente sua rede agora mesmo!
            </Typography>
          </Box>
          <Typography variant={"body1"}>
            ⭐ Você tem {userStatistics.totalFans} fãs.
          </Typography>
          <Typography variant={"body1"}>
            📝 {userStatistics.totalScraps} recados.
          </Typography>
          <Typography variant={"body1"}>
            <strong>Visualizações do seu perfil:</strong>
          </Typography>
          <Typography variant={"body1"}>
            Desde {dayjsFromDate(user.createdAt, "D [de] MMMM [de] YYYY")}:{" "}
            {userStatistics.totalProfileViews}. Semana passada:{" "}
            {userStatistics.lastWeekProfileViews}.
          </Typography>
          <Box>
            <Typography variant={"body1"}>
              <strong>Visitantes recentes:</strong>
            </Typography>
            <Typography variant={"body1"}>
              {userStatistics.profileLastVisitors
                .map(({ name }) => name)
                .join(", ")}
            </Typography>
          </Box>
          <Box>
            <Typography variant={"body1"}>
              <strong>Próximas etapas:</strong>
            </Typography>
            <Typography variant={"body1"}>Classifique seus amigos!</Typography>
          </Box>
          <Box>
            <Typography variant={"body1"}>
              <strong>Sorte de hoje</strong>
            </Typography>
            <Typography variant={"body1"}>{todaysLuck}</Typography>
          </Box>
        </Stack>
      </Paper>
    </Box>
  );
};
