import { Skeleton, Stack, Typography } from "@mui/material";
import { useCommunityHotTopics } from "api/community/queries";
import { useMeAmMember } from "api/me/queries";
import { Button } from "components/Button";
import { BoxWithTitle } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { MembersOnly } from "features/community/components/MembersOnly";
import { useForumAccess } from "features/community/useForumAccess";

interface IProps {
  communityId: number;
}
export const HotTopics: React.FC<IProps> = ({ communityId }) => {
  const { canRead, isLoading: isLoadingAccess } = useForumAccess(communityId);
  const { data: topics, isLoading: isLoadingData, isFetching } = useCommunityHotTopics(communityId, canRead);
  // The request waits for the access check; both count as loading.
  const isLoading = isLoadingAccess || isLoadingData;
  const { data: membership } = useMeAmMember(communityId);

  const navigate = useNavigate();

  const handleClickSeeAll = () => {
    navigate(`forum`)
  }

  const handleCreateTopic = () => {
    navigate(`forum/new-topic`)
  }

  const items =
    topics?.topics?.map((t) => (
      <Stack direction={"row"} gap={1} alignItems={"center"}>
        {/* flex, so the inline picture doesn't get a text-line gap under it. */}
        <Link to={`/profile/${t.authorId}`} style={{ display: "flex" }}>
          <img src={t.authorProfilePicture ?? undefined} height={40} />
        </Link>
        <Stack direction={"column"} sx={{ overflow: "hidden", minWidth: 0 }}>
          <Link to={`forum/topic/${t.id}`}><Typography>{t.title}</Typography></Link>
          <Typography variant={"body2"}>{t.replies} respostas</Typography>
          {/* Topics without replies have no last reply (the API sends a zero date). */}
          {t.replies > 0 && (
            <Stack
              direction={"row"}
              alignItems={"center"}
              spacing={1}
              sx={{
                maxWidth: "100%",
                whiteSpace: "nowrap",
              }}>
              <Typography variant={"body2"} sx={{ flexShrink: 0 }}>
                {dayjsFromDate(t.lastReplyDate, "DD [de] MMM [de] YYYY")}:
              </Typography>
              {/* One line, cut with "…" where the box ends. */}
              <Typography variant={"body2"} noWrap sx={{ minWidth: 0 }}>{t.lastReplyPreview}</Typography>
            </Stack>
          )}
        </Stack>
      </Stack>
    )) ?? [];

  if (!isLoadingAccess && !canRead) {
    return (
      <BoxWithTitle title={"fórum"}>
        <MembersOnly communityId={communityId} action={"ver o fórum desta comunidade privada"} />
      </BoxWithTitle>
    );
  }

  return (
    <BoxWithTitle title={"fórum"} showProgress={isFetching && !isLoading}>
      <StripedGrid
        items={items}
        elementsPerRow={1}
        isLoading={isLoading}
        skeleton={
          <Stack direction={"row"} gap={1} alignItems={"center"}>
            <Skeleton variant={"rectangular"} width={40} height={40} sx={{ minWidth: 40 }} />
            <Stack direction={"column"}>
              <Typography><Skeleton variant={"text"} width={160} /></Typography>
              <Typography variant={"body2"}><Skeleton variant={"text"} width={80} /></Typography>
            </Stack>
          </Stack>
        }
      />
      <Stack direction={"row"} justifyContent={"end"}>
        <Button onClick={handleClickSeeAll}>ver todos</Button>
        {membership?.amIMember && <Button onClick={handleCreateTopic}>+ novo tópico</Button>}
      </Stack>
    </BoxWithTitle>
  )
}
