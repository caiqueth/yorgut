import React from "react";

import { Box, LinearProgress, Paper, Stack, Typography } from "@mui/material";
import { InlineSkeleton } from "components/InlineSkeleton";
import { dayjsFromDate } from "helpers/formatting/date";
import { useMeWelcome } from "api/me/queries";
import { Link } from "react-router-dom";
import { WelcomeResponse } from "api/generated/model";

// The "próximas etapas" that still apply: find friends, rate them (only once
// there are friends to rate), join a community. Empty when all are done.
const nextSteps = (welcome: WelcomeResponse) => [
  ...(welcome.directConnections === 0 ? [{ text: "Encontre amigos!", to: "/search" }] : []),
  ...(welcome.directConnections > 0 && !welcome.hasRatedFriends
    ? [{ text: "Classifique seus amigos!", to: "/friends" }]
    : []),
  ...(welcome.communityCount === 0
    ? [{ text: "Encontre uma comunidade!", to: "/communities/discover" }]
    : []),
];

// A number of the welcome text, or a skeleton of about its width until it loads.
const value = (n?: number) => n ?? <InlineSkeleton width={"2ch"} />;

export const Welcome: React.FC = () => {
  const { data, isFetching, isLoading } = useMeWelcome();

  return (
    <Box>
      <Paper sx={{ p: 1, borderRadius: (theme) => theme.yorgut.featuredBoxRadius }}>
        {/* Reloads (after a scrap or a friendship change); hidden, not removed, so nothing moves. */}
        <LinearProgress sx={{ visibility: isFetching && !isLoading ? "visible" : "hidden", mb: 0.5 }} />
        <Stack direction={"column"} gap={1}>
          <Typography
            variant={"subtitle1"}
            sx={{ "&::first-letter": { fontSize: 36, lineHeight: 1 } }}
          >
            Bem-vindo(a), {data ? data.userFirstName : <InlineSkeleton width={"6ch"} />}
          </Typography>
          <Box>
            <Typography variant={"body1"}>
              Você está conectado(a) a {value(data?.indirectConnections)}{" "}
              pessoas através de {value(data?.directConnections)} amigos
            </Typography>
            <Typography variant={"body1"}>
              Aumente sua rede agora mesmo!
            </Typography>
          </Box>
          <Typography variant={"body1"}>
            ⭐ Você tem {value(data?.numberOfFans)} fãs.
          </Typography>
          <Typography variant={"body1"}>
            📝 {value(data?.numberOfScraps)} recados.
          </Typography>
          <Typography variant={"body1"}>
            <strong>Visualizações do seu perfil:</strong>
          </Typography>
          <Typography variant={"body1"}>
            Desde {data
              ? dayjsFromDate(data.userCreatedAt, "D [de] MMMM [de] YYYY")
              : <InlineSkeleton width={"18ch"} />}:{" "}
            {value(data?.profileViews.allTime)}. Últimos 7 dias:{" "}
            {value(data?.profileViews.lastSevenDays)}.
          </Typography>
          <Box>
            <Typography variant={"body1"}>
              <strong>Visitantes recentes:</strong>
            </Typography>
            <Typography variant={"body1"}>
              {!data
                ? <InlineSkeleton width={"24ch"} />
                : !data.recordVisits
                ? <>Você desligou o registro de visitas em <Link to="/my-account">minha conta</Link>.</>
                : data.recentVisitors.length > 0
                ? data.recentVisitors.map(({ userId, firstName }, idx) => (
                  <React.Fragment key={userId}>
                    {idx > 0 && ", "}
                    <Link to={`/profile/${userId}`}>{firstName}</Link>
                  </React.Fragment>
                ))
                : "Por enquanto, ninguém."}
            </Typography>
          </Box>
          {/* Only the steps still to do, so nothing to show until the data says which. */}
          {data && nextSteps(data).length > 0 && (
            <Box>
              <Typography variant={"body1"}>
                <strong>Próximas etapas:</strong>
              </Typography>
              {nextSteps(data).map((step) => (
                <Typography key={step.to} variant={"body1"}>
                  <Link to={step.to}>{step.text}</Link>
                </Typography>
              ))}
            </Box>
          )}
          <Box>
            <Typography variant={"body1"}>
              <strong>Sorte de hoje</strong>
            </Typography>
            <Typography variant={"body1"}>
              {data ? data.todaysLuck : <InlineSkeleton width={"36ch"} />}
            </Typography>
          </Box>
        </Stack>
      </Paper>
    </Box>
  );
};
