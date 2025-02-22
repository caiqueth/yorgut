import { useTodaysLuck, useUser, useUserStatistics } from "api/user/queries";

import { Box, Paper, Stack, Typography } from "@mui/material";

export const Welcome = () => {
  const userId = 1;
  const { data: user } = useUser(userId);
  const userStatistics = useUserStatistics(userId);
  const todaysLuck = useTodaysLuck(userId);

  return (
    <Box>
      <Paper sx={{ p: 1 }}>
        <Stack direction={"column"} gap={1}>
          <Typography variant={"subtitle1"} sx={{ "&::first-letter": { fontSize: "200%" } }}>
            Bem-vindo(a), {user?.firstName}
          </Typography>
          <Box>
            <Typography variant={"body1"}>
              Você está conectado a {userStatistics.undirectConnections} pessoas através de{" "}
              {userStatistics.totalFriends} amigos
            </Typography>
            <Typography variant={"body1"}>Aumente sua rede agora mesmo!</Typography>
          </Box>
          <Typography variant={"body1"}>⭐ Você tem {userStatistics.totalFans} fãs.</Typography>
          <Typography variant={"body1"}>📝 {userStatistics.totalScraps} recados.</Typography>
          <Typography variant={"body1"}>
            <strong>Visualizações do seu perfil:</strong>
          </Typography>
          <Typography variant={"body1"}>
            Desde {user?.createdAt}: {userStatistics.totalProfileViews}, Semana passada:{" "}
            {userStatistics.lastWeekProfileViews}
          </Typography>
          <Box>
            <Typography variant={"body1"}>
              <strong>Visitantes recentes:</strong>
            </Typography>
            <Typography variant={"body1"}>
              {userStatistics.profileLastVisitors.map(({ name }) => name).join(", ")}
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
