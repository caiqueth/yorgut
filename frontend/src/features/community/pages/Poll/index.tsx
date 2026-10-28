import React, { useState } from "react";
import { Box, FormControlLabel, LinearProgress, Radio, RadioGroup, Skeleton, Stack, Typography } from "@mui/material";
import { Link, useNavigate } from "react-router-dom";
import { useClosePoll, usePoll, useVotePoll } from "api/poll/queries";
import { PollResponse } from "api/generated/model";
import { BoxWithTitle } from "components/BoxWithTitle";
import { Button } from "components/Button";
import { MembersOnly } from "features/community/components/MembersOnly";
import { useForumAccess } from "features/community/useForumAccess";
import { pollStatusText } from "features/community/pollStatus";
import { dayjsFromDate } from "helpers/formatting/date";
import { useIdParam } from "hooks/useIdParam";

// Each option with a bar of its share of the votes; the logged user's choice
// is marked.
const PollResults: React.FC<{ poll: PollResponse }> = ({ poll }) => (
  <Stack gap={1}>
    {poll.options.map((o) => {
      const votes = o.votes ?? 0;
      const percent = poll.voteCount ? Math.round((votes / poll.voteCount) * 100) : 0;
      return (
        <Box key={o.id}>
          <Typography variant={"body2"}>
            {o.text}{o.id === poll.myOptionId && " ✓"}
            {" — "}{percent}% ({votes} {votes === 1 ? "voto" : "votos"})
          </Typography>
          <LinearProgress variant={"determinate"} value={percent} sx={{ height: 8 }} />
        </Box>
      );
    })}
  </Stack>
);

const PollVote: React.FC<{ poll: PollResponse }> = ({ poll }) => {
  const [optionId, setOptionId] = useState<number | null>(null);
  const { mutate: vote, isPending } = useVotePoll(poll.id, poll.communityId);

  return (
    <Box>
      <RadioGroup value={optionId ?? ""} onChange={(e) => setOptionId(Number(e.target.value))}>
        {poll.options.map((o) => (
          <FormControlLabel key={o.id} value={o.id} control={<Radio size={"small"} />} label={o.text} />
        ))}
      </RadioGroup>
      <Button onClick={() => optionId !== null && vote(optionId)} disabled={optionId === null} loading={isPending}>
        votar
      </Button>
    </Box>
  );
};

// A poll: members who haven't voted pick an option; the results show once
// you voted or the poll closed.
export const Poll: React.FC = () => {
  const communityId = useIdParam("communityId");
  const pollId = useIdParam("pollId");
  const { canRead, isLoading: isLoadingAccess } = useForumAccess(communityId);
  const { data: poll, isFetching } = usePoll(pollId, canRead);
  const { mutate: close, isPending: isClosing } = useClosePoll(pollId, communityId);
  const navigate = useNavigate();

  if (!isLoadingAccess && !canRead) {
    return (
      <BoxWithTitle title={"enquete"}>
        <MembersOnly communityId={communityId} action={"ver as enquetes desta comunidade privada"} />
      </BoxWithTitle>
    );
  }

  if (!poll) {
    return (
      <BoxWithTitle title={"enquete"} showProgress>
        <Box sx={{ p: 1 }}>
          <Skeleton variant={"text"} width={"60%"} />
          <Skeleton variant={"text"} width={"40%"} />
        </Box>
      </BoxWithTitle>
    );
  }

  const showResults = poll.options.some((o) => o.votes !== null);

  return (
    <BoxWithTitle title={"enquete"} showProgress={isFetching}>
      <Stack sx={{ p: 1 }} gap={1}>
        <Box>
          <Typography variant={"h6"}>{poll.question}</Typography>
          <Typography variant={"body2"}>
            por <Link to={`/profile/${poll.authorId}`}>{poll.authorFirstName}</Link>
            {" em "}{dayjsFromDate(poll.createdAt, "D [de] MMMM [de] YYYY")}
            {" · "}{pollStatusText(poll)}
            {" · "}{poll.voteCount} {poll.voteCount === 1 ? "voto" : "votos"}
          </Typography>
        </Box>
        {poll.description && <Typography sx={{ whiteSpace: "pre-wrap" }}>{poll.description}</Typography>}

        {poll.canVote ? (
          <PollVote poll={poll} />
        ) : showResults ? (
          <PollResults poll={poll} />
        ) : (
          <Box>
            {poll.options.map((o) => (
              <Typography key={o.id}>• {o.text}</Typography>
            ))}
            <Typography variant={"body2"} sx={{ mt: 1 }}>
              Só membros da comunidade votam. Os resultados aparecem quando a enquete for encerrada.
            </Typography>
          </Box>
        )}

        <Stack direction={"row"} gap={1}>
          <Button onClick={() => navigate(`/community/${communityId}/polls`)}>voltar às enquetes</Button>
          {poll.canClose && (
            <Button onClick={() => close()} loading={isClosing}>encerrar enquete</Button>
          )}
        </Stack>
      </Stack>
    </BoxWithTitle>
  );
};
